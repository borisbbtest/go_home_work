package storage

type StorageURL struct {
	Port int
	Url  string
	Path string
}

type repositoriesURLShort interface {
	getURLforRedirect(urlshotr string)
}

func (store *StorageURL) setURLforRedirect() (err error, dataStore StorageURL) {

	return
}

func (store *StorageURL) getURLforRedirect(urlshort string) {

}
