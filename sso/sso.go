package sso

import (
	"errors"
	tokenGo "github.com/weloe/token-go"
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
)

// Options construct options
type Options struct {
	SsoOptions  *config.SsoOptions
	SignOptions *config.SignOptions
	Enforcer    tokenGo.IEnforcer
}

type SsoEnforcer struct {
	apiName    *ApiName
	paramName  *ParamName
	config     *config.SsoConfig
	signConfig *config.SignConfig
	enforcer   tokenGo.IEnforcer
}

// NewSsoEnforcer create sso enforcer.
// If the available field in the parameter is empty, use the default value,
// if the required field is empty, return nil and error.
func NewSsoEnforcer(options *Options) (*SsoEnforcer, error) {
	if options.Enforcer == nil {
		return nil, errors.New("Options.Enforcer can not be nil")
	}
	if options.SsoOptions.CookieDomain != "" {
		options.Enforcer.GetTokenConfig().CookieConfig.Domain = options.SsoOptions.CookieDomain
	}

	ssoConfig, err := config.NewSsoConfig(options.SsoOptions)
	if err != nil {
		return nil, err
	}
	signConfig, err := config.NewSignConfig(options.SignOptions)
	if err != nil {
		return nil, err
	}

	return &SsoEnforcer{
		apiName:    DefaultApiName(),
		paramName:  DefaultParamName(),
		config:     ssoConfig,
		signConfig: signConfig,
		enforcer:   options.Enforcer,
	}, nil
}

func (s *SsoEnforcer) SetApi(apiName *ApiName) {
	s.apiName = apiName
}

func (s *SsoEnforcer) GetApi() *ApiName {
	return s.apiName
}

func (s *SsoEnforcer) SetParamName(paramName *ParamName) {
	s.paramName = paramName
}

func (s *SsoEnforcer) GetParamName() *ParamName {
	return s.paramName
}

func (s *SsoEnforcer) SetSsoConfig(config *config.SsoConfig) {
	s.config = config
}

func (s *SsoEnforcer) GetSsoConfig() *config.SsoConfig {
	return s.config
}

func (s *SsoEnforcer) SetSignConfig(config *config.SignConfig) {
	s.signConfig = config
}

func (s *SsoEnforcer) GetSignConfig() *config.SignConfig {
	return s.signConfig
}

// ssoLogoutBack single-logout callback, redirect to ParamName.Back url.
// If http request has back param and value is SELF, redirect to previous page,
// if http request back param and value is validated url, redirect to url,
// else return model.Ok() directly.
func (s *SsoEnforcer) ssoLogoutBack(ctx ctx.Context) (interface{}, error) {
	paramName := s.paramName
	request := ctx.Request()
	response := ctx.Response()
	back := request.Query(paramName.Back)
	if back != "" {
		if back == constant.SELF {
			return "<script>if(document.referrer != location.href){ location.replace(document.referrer || '/'); }</script>", nil
		}
		response.Redirect(back)
		return nil, nil
	} else {
		// back is nil
		return model.Ok().SetMsg("back is nil, not redirect"), nil
	}

}
