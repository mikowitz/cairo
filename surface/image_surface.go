package surface

import "github.com/mikowitz/cairo/status"

type ImageSurface struct {
	*BaseSurface
	format        Format
	width, height int
}

func NewImageSurface(format Format, width, height int) (*ImageSurface, error) {
	ptr := imageSurfaceCreate(format, width, height)
	st := surfaceStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	baseSurface := newBaseSurface(ptr)
	return &ImageSurface{
		BaseSurface: baseSurface,
		format:      format,
		width:       width,
		height:      height,
	}, nil
}

func (s *ImageSurface) GetFormat() Format {
	s.RLock()
	defer s.RUnlock()

	return s.format
}

func (s *ImageSurface) GetWidth() int {
	s.RLock()
	defer s.RUnlock()

	return s.width
}

func (s *ImageSurface) GetHeight() int {
	s.RLock()
	defer s.RUnlock()

	return s.height
}

func (s *ImageSurface) GetStride() int {
	s.RLock()
	defer s.RUnlock()

	return s.format.StrideForWidth(s.width)
}
