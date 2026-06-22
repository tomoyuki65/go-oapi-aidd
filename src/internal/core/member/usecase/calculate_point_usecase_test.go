//go:build unit

package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-oapi-aidd/internal/core/member/domain"
	mockRepository "go-oapi-aidd/internal/core/member/domain/mock_repository"
)

func TestCalculatePointUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mockRepository.NewMockMemberRepository(ctrl)

	t.Run("会員取得成功時に正しい計算結果を返すこと", func(t *testing.T) {
		rank, err := domain.NewRank("silver")
		assert.NoError(t, err)
		member, err := domain.NewMember("22222222-2222-2222-2222-222222222222", "佐藤", rank)
		assert.NoError(t, err)
		repository.EXPECT().
			FindByID(gomock.Any(), gomock.Nil(), "22222222-2222-2222-2222-222222222222").
			Return(member, nil)

		usecase := NewCalculatePointUsecase(repository)
		output, err := usecase.Execute(context.Background(), CalculatePointInput{
			MemberID:       "22222222-2222-2222-2222-222222222222",
			PurchaseAmount: 5000,
		})

		assert.NoError(t, err)
		assert.Equal(t, "22222222-2222-2222-2222-222222222222", output.MemberID)
		assert.Equal(t, 5000, output.PurchaseAmount)
		assert.Equal(t, 150, output.GrantedPoint)
	})

	t.Run("repositoryの会員未存在エラーをMEMBER_NOT_FOUND相当へ扱うこと", func(t *testing.T) {
		repository.EXPECT().
			FindByID(gomock.Any(), gomock.Nil(), "99999999-9999-9999-9999-999999999999").
			Return(nil, domain.ErrMemberNotFound)

		usecase := NewCalculatePointUsecase(repository)
		_, err := usecase.Execute(context.Background(), CalculatePointInput{
			MemberID:       "99999999-9999-9999-9999-999999999999",
			PurchaseAmount: 5000,
		})

		assert.ErrorIs(t, err, ErrMemberNotFound)
	})

	t.Run("repositoryの予期しないエラーを内部エラー相当へ扱うこと", func(t *testing.T) {
		repository.EXPECT().
			FindByID(gomock.Any(), gomock.Nil(), "11111111-1111-1111-1111-111111111111").
			Return(nil, assert.AnError)

		usecase := NewCalculatePointUsecase(repository)
		_, err := usecase.Execute(context.Background(), CalculatePointInput{
			MemberID:       "11111111-1111-1111-1111-111111111111",
			PurchaseAmount: 5000,
		})

		assert.ErrorIs(t, err, ErrInternal)
	})
}
