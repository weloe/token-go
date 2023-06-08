package config

// SignConfig sign config
type SignConfig struct {
	SecretKey          string
	TimeStampDisparity int64
	IsCheckNonce       bool
}

func NewSignConfig(options *SignOptions) (*SignConfig, error) {
	if options == nil {
		options = &SignOptions{}
	}
	if options.TimeStampDisparity == 0 {
		options.TimeStampDisparity = 1000 * 60 * 15
	}
	return &SignConfig{
		SecretKey:          options.SecretKey,
		TimeStampDisparity: options.TimeStampDisparity,
		IsCheckNonce:       options.IsCheckNonce,
	}, nil
}

func (s *SignConfig) GetSaveNonceExpire() int64 {
	if s.TimeStampDisparity >= 0 {
		return s.TimeStampDisparity / 1000
	} else {
		return 60 * 60 * 24
	}
}
