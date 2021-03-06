package server

import (
	"crypto/tls"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/json"

	"net/http"
)

var (
	scheme          = runtime.NewScheme()
	codecs          = serializer.NewCodecFactory(scheme)
	tlscert, tlskey string
)

type AdmissionController interface {
	HandleAdmission(review *v1beta1.AdmissionReview) error
}

type AdmissionControllerServer struct {
	AdmissionController AdmissionController
	Decoder             runtime.Decoder
}

func (acs *AdmissionControllerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if data, err := ioutil.ReadAll(r.Body); err == nil {
		body = data
	}
	logrus.Debugln(string(body))
	review := &v1beta1.AdmissionReview{}
	_, _, err := acs.Decoder.Decode(body, nil, review)
	if err != nil {
		logrus.Errorln("Can't decode request", err)
	}
	acs.AdmissionController.HandleAdmission(review)
	responseInBytes, err := json.Marshal(review)

	if _, err := w.Write(responseInBytes); err != nil {
		logrus.Errorln(err)
	}
}

func GetAdmissionServerNoSSL(ac AdmissionController, listenOn string) *http.Server {
	server := &http.Server{
		Handler: &AdmissionControllerServer{
			AdmissionController: ac,
			Decoder:             codecs.UniversalDeserializer(),
		},
		Addr: listenOn,
	}

	return server
}

func GetAdmissionValidationServer(ac AdmissionController, tlsCert, tlsKey, listenOn string) *http.Server {
	sCert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	server := GetAdmissionServerNoSSL(ac, listenOn)
	server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
	if err != nil {
		logrus.Error(err)
	}
	return server
}
