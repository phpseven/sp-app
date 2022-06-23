package bdetect

import (
    "encoding/json"
    "fmt"
    "github.com/avct/uasurfer"
    rhttp "github.com/spiral/roadrunner/service/http"
    "github.com/spiral/roadrunner/service/http/attributes"
    "net/http"
)

const ID = "bdetect"

type Service struct{}

func (s *Service) Init(rhttp *rhttp.Service) (bool, error) {
    rhttp.AddMiddleware(s.middleware)

    return true, nil
}

func (s *Service) middleware(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ua := uasurfer.Parse(r.Header.Get("User-Agent"))
        data, _ := json.Marshal(struct {
            Browser        string `json:"browser"`
            BrowserVersion string `json:"browserVersion"`
            IsBot          bool   `json:"isBot"`
        }{
            Browser: ua.Browser.Name.StringTrimPrefix(),
            BrowserVersion: fmt.Sprintf(
                "%v.%v.%v",
                ua.Browser.Version.Major,
                ua.Browser.Version.Minor,
                ua.Browser.Version.Patch,
            ),
            IsBot: ua.IsBot(),
        })

        attributes.Set(r, "bdetect", string(data))

        f(w, r)
    }
}
