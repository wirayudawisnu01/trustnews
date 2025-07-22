package service

import (
	"context"
	"trustnews/config"
	"trustnews/internal/adapter/cloudflare"
	"trustnews/internal/adapter/repository"
	"trustnews/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
	UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	cfg         *config.Config
	r2          cloudflare.CloudflareR2Adapter
}

// CreateContent implements ContentService.
func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepo.CreateContent(ctx, req)
	if err != nil  {
		code = "[SERVICE] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteContent implements ContentService.
func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	err = c.contentRepo.DeleteContent(ctx, id)
	if err != nil {
		code = "[SERVICE] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetContentByID implements ContentService.
func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	result, err := c.contentRepo.GetContentByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

// GetContents implements ContentService.
func (c *contentService) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error) {
	results, totalData, totalPages, err := c.contentRepo.GetContents(ctx, query)
	if err != nil {
		code = "[SERVICE] GetContents - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPages, nil
}

// UpdateContent implements ContentService.
func (c *contentService) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepo.UpdateContent(ctx, req)
	if err != nil {
		code = "[SERVICE] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// UploadImageR2 implements ContentService.
func (c *contentService) UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error) {
	urlImage, err := c.r2.UploadImage(&req)
	if err != nil {
		code = "[SERVICE] UploadImageR2 - 1"
		log.Errorw(code, err)
		return "", err
	}

	return urlImage, nil
}

func NewContentService(repo repository.ContentRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) ContentService {
	return &contentService{
		contentRepo: repo,
		cfg:         cfg,
		r2:          r2,
	}
}
