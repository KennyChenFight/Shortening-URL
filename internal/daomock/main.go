package daomock

//go:generate mockgen -destination=mock.go -package=$GOPACKAGE github.com/KennyChenFight/Shortening-URL/pkg/dao CacheDAO,KeyDAO,UrlDAO
