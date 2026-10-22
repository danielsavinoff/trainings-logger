// START_MODULE_CONTRACT
//   PURPOSE: Persist only explicitly completed questionnaire answers.
//   SCOPE: Validate optional preferences and authorized accounts, then save atomically.
//   DEPENDS: M-DOMAIN-USER-PREFERENCES, M-DOMAIN-IDENTITY
//   LINKS: M-APP-QUESTIONNAIRE, V-M-APP-QUESTIONNAIRE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   CompleteQuestionnaire - carry optional completed answers
//   CompleteQuestionnaireHandler - retain repositories and transaction boundary
//   NewCompleteQuestionnaireHandler - wire answer persistence
//   Handle - validate and atomically persist explicit answers
// END_MODULE_MAP

package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/identity"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

// START_BLOCK_QUESTIONNAIRE_COMPLETION
type CompleteQuestionnaire struct {
	WeightUnit     *shared.WeightUnit
	LinkedIdentity *ports.LinkedIdentity
}

type CompleteQuestionnaireHandler struct {
	preferencesRepository preferences.PreferencesRepository
	identityRepository    identity.IdentityRepository
	credentialsRepository identity.IdentityCredentialsRepository
	transactions          ports.TransactionRunner[ports.TransactionNoResult]
}

func NewCompleteQuestionnaireHandler(
	preferencesRepository preferences.PreferencesRepository,
	identityRepository identity.IdentityRepository,
	credentialsRepository identity.IdentityCredentialsRepository,
	transactions ports.TransactionRunner[ports.TransactionNoResult],
) CompleteQuestionnaireHandler {
	return CompleteQuestionnaireHandler{
		preferencesRepository: preferencesRepository,
		identityRepository:    identityRepository,
		credentialsRepository: credentialsRepository,
		transactions:          transactions,
	}
}

func (h CompleteQuestionnaireHandler) Handle(
	ctx context.Context,
	input CompleteQuestionnaire,
) error {
	if input.WeightUnit == nil && input.LinkedIdentity == nil {
		return nil
	}

	var prefs preferences.Preferences
	if input.WeightUnit != nil {
		validated, err := preferences.NewPreferences(*input.WeightUnit)
		if err != nil {
			return err
		}
		prefs = validated
	}

	var account identity.Identity
	var credentials identity.IdentityCredentials
	if input.LinkedIdentity != nil {
		var err error
		account, err = identity.NewIdentity(
			input.LinkedIdentity.Identity.ID,
			input.LinkedIdentity.Identity.Provider,
		)
		if err != nil {
			return err
		}

		value := input.LinkedIdentity.Credentials
		credentials, err = identity.NewIdentityCredentials(
			value.IdentityID,
			value.AccessTokenEncrypted,
			value.RefreshTokenEncrypted,
			value.AccessTokenExpiresAt,
			value.UpdatedAt,
		)
		if err != nil {
			return err
		}
		if credentials.IdentityID != account.ID {
			return identity.ErrInvalidCredentials
		}
	}

	_, err := h.transactions.Run(
		ctx,
		func(ctx context.Context) (ports.TransactionNoResult, error) {
			if input.WeightUnit != nil {
				if err := h.preferencesRepository.Save(
					ctx,
					prefs,
				); err != nil {
					return ports.TransactionNoResult{}, err
				}
			}

			if input.LinkedIdentity != nil {
				if err := h.identityRepository.Save(
					ctx,
					account,
				); err != nil {
					return ports.TransactionNoResult{}, err
				}
				if err := h.credentialsRepository.Save(
					ctx,
					credentials,
				); err != nil {
					return ports.TransactionNoResult{}, err
				}
			}

			return ports.TransactionNoResult{}, nil
		},
	)

	return err
}

// END_BLOCK_QUESTIONNAIRE_COMPLETION
