package providers

import "pay-system/ports"

type ThirdParty struct {
	provider ports.IThirdPartyService
}

func NewThirdParty(provider ports.IThirdPartyService) *ThirdParty {
	return &ThirdParty{
		provider: provider,
	}
}
