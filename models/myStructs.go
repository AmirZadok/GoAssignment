package models

type InfoMessage struct {
	Message string
}

type Host struct {
	Id                     int
	Uuid, Name, Ip_Address string
}

type Container struct {
	Id, Host_Id      int
	Name, Image_Name string
}

type PartialContainer struct {
	Host_Id    int
	Image_Name string
}

type ContainerWithHostName struct {
	Id, Host_Id      int
	Name, Image_Name , Host_Name string
}