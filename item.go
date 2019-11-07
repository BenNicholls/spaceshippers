package main

//Items are things that can be stored in the storage system
type Item struct {
	Name        string
	Description string
	Volume      float64     // (L)
	StorageType StorageType //which type of storage is used ex. General, liquid, gas, etc.
}

func (i Item) GetName() string {
	return i.Name
}

func (i Item) GetAmount() float64 {
	return i.Volume
}

func (i Item) GetStorageType() StorageType {
	return i.StorageType
}

func (i *Item) ChangeAmount(d float64) {
	i.Volume += d
	if i.Volume < 0 {
		i.Volume = 0
	}
}

func (i *Item) SetAmount(v float64) {
	i.Volume = v
}

func (i *Item) GetDescription() string {
	if i.Description != "" {
		return i.Description
	}

	return "This item has no description. It is non-descript."
}
