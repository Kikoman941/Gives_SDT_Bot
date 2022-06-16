package images

type Repository interface {
	SaveImage(img string) error
}
