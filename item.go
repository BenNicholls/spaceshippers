package main

//Items are things can be stored in the storage system
type Item struct {
	Name        string
	Volume      float64  // (L)
	StorageType StatType //which type of storage is used ex. Cold, volatile, etc
}

func (i Item) GetName() string {
	return i.Name
}

func (i Item) GetVolume() float64 {
	return i.Volume
}

func (i Item) GetStorageType() StatType {
	return i.StorageType
}

func (i *Item) ChangeVolume(d float64) {
	i.Volume += d
	if i.Volume < 0 {
		i.Volume = 0
	}
}

func (i *Item) SetVolume(v float64) {
	i.Volume = v
}
