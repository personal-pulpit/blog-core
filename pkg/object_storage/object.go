package objStorage

type Object struct {
	Name        string
	Data        []byte
	Size        int64
	ContentType string
}

func NewObject(name string, data []byte, size int64, contentType string) *Object {
	return &Object{
		Name:        name,
		Data:        data,
		Size:        size,
		ContentType: contentType,
	}
}

const ImageContentType = "image/webp"
const ProfileImageBucketName = "blog"