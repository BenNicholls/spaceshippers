package main

import "errors"

type Storable interface {
	GetName() string
	GetVolume() float64
	GetStorageType() StatType //determins which kind fo storage needs to be used
	ChangeVolume(d float64)
	SetVolume(v float64)
}

type StorageSystem struct {
	SystemStats

	items []Storable

	volume map[StatType]float64

	ship *Ship
}

func NewStorageSystem(s *Ship) *StorageSystem {
	ss := new(StorageSystem)

	ss.ship = s

	ss.volume = make(map[StatType]float64)
	ss.volume[STAT_GENERAL_STORAGE] = 0
	ss.volume[STAT_VOLATILE_STORAGE] = 0
	ss.volume[STAT_FUEL_STORAGE] = 0
	ss.volume[STAT_COLD_STORAGE] = 0

	return ss
}

func (ss *StorageSystem) Update(tick int) {
	//put love here
}

func (ss *StorageSystem) AddToStores(item Storable) error {

	if float64(ss.SystemStats.GetStat(item.GetStorageType()))-ss.volume[item.GetStorageType()] < item.GetVolume() {
		return errors.New("Not enough space")
	} else {
		ss.volume[item.GetStorageType()] += item.GetVolume()
	}

	for _, i := range ss.items {
		if i.GetName() == item.GetName() {
			i.ChangeVolume(item.GetVolume())
			return nil
		}
	}

	ss.items = append(ss.items, item)

	return nil
}

func (ss *StorageSystem) RemoveFromStores(name string, v float64) (item Storable, err error) {
	for n, i := range ss.items {
		if i.GetName() == name { //item found
			if i.GetVolume() >= v { //enough is there
				item = &Item{
					Name:        i.GetName(),
					Volume:      v,
					StorageType: i.GetStorageType(),
				}
				i.ChangeVolume(-item.GetVolume())
			} else { //not enough there. return what we got
				item = &Item{
					Name:        i.GetName(),
					Volume:      i.GetVolume(),
					StorageType: i.GetStorageType(),
				}
				ss.items = append(ss.items[:n], ss.items[n+1:]...)
				err = errors.New("Insufficient amount in stores")

			}
		}

		ss.volume[item.GetStorageType()] -= item.GetVolume()
		err = errors.New("No item found")
	}

	return
}

func (ss *StorageSystem) GetVolume(name string) float64 {
	for _, i := range ss.items {
		if i.GetName() == name {
			return i.GetVolume()
		}
	}

	return 0
}
