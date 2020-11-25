package main

import (
	"errors"

	"github.com/bennicholls/burl-E/burl"
)

type Storable interface {
	GetName() string
	GetDescription() string
	GetAmount() float64 //returns amount of item in (L) for items and liquid, or in molar value for gasses
	ChangeAmount(d float64)
	SetAmount(v float64)
	GetStorageType() StorageType //which kind of storage needs to be used
}

type StorageSystem struct {
	SystemStats

	items map[string]Storable

	volume   [STORE_MAXTYPES]float64 //General and liquid storage volume/capacity is in litres, gas is in molar value (pressure*volume)
	capacity [STORE_MAXTYPES]float64

	ship *Ship
}

type StorageType int

const (
	STORE_GENERAL StorageType = iota
	STORE_LIQUID
	STORE_GAS

	STORE_MAXTYPES
)

func NewStorageSystem(s *Ship) *StorageSystem {
	ss := new(StorageSystem)
	ss.items = make(map[string]Storable)
	ss.ship = s

	return ss
}

func (ss *StorageSystem) Update(tick int) {
	//put love here
}

//ensure storage capacities are updated if stats change
func (ss *StorageSystem) OnStatUpdate() {
	ss.capacity[STORE_GENERAL] = float64(ss.GetStat(STAT_GENERAL_STORAGE))
	ss.capacity[STORE_LIQUID] = float64(ss.GetStat(STAT_LIQUID_STORAGE))
	ss.capacity[STORE_GAS] = float64(ss.GetStat(STAT_GAS_STORAGE) * 50000) //NOTE: currently limiting gas storage to 50000 kPa
	burl.LogInfo("updated capacties")
	burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "stores"))
}

func (ss *StorageSystem) Store(item Storable) error {
	if ss.capacity[item.GetStorageType()]-ss.volume[item.GetStorageType()] < item.GetAmount() {
		return errors.New("Not enough space")
	}

	if item.GetAmount() == 0 {
		return errors.New("Cannot store zero of item.")
	}

	ss.volume[item.GetStorageType()] += item.GetAmount()

	if i, ok := ss.items[item.GetName()]; ok {
		i.ChangeAmount(item.GetAmount())
	} else {
		ss.items[item.GetName()] = item
	}

	burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "stores"))

	return nil
}

//Attempts to remove item from stores. Returns the amount of item removed, or returns 0 if no item is found.
//If less than volume v is present in stores, returns just what was there and removes the item from the
//ship's inventory entirely. Check err to see what the deal is.
func (ss *StorageSystem) Remove(item Storable) (amount float64, err error) {
	if i, ok := ss.items[item.GetName()]; !ok {
		return 0, errors.New("Item not found.")
	} else if i.GetAmount() <= item.GetAmount() {
		if i.GetAmount() != item.GetAmount() {
			err = errors.New("Insufficient amount of item in stores")
		}
		delete(ss.items, i.GetName())
		burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "stores"))
		return i.GetAmount(), err
	} else {
		i.ChangeAmount(-item.GetAmount())
		burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "stores"))
		return item.GetAmount(), nil
	}
}

func (ss *StorageSystem) GetItemVolume(name string) float64 {
	if i, ok := ss.items[name]; !ok {
		return 0
	} else {
		return i.GetAmount()
	}
}

func (ss *StorageSystem) GetFilledVolume(storageType StorageType) float64 {
	return ss.volume[storageType]
}

func (ss *StorageSystem) GetCapacity(storageType StorageType) float64 {
	return ss.capacity[storageType]
}

func (ss *StorageSystem) GetFillPct(st StorageType) float64 {
	return 100 * ss.volume[st] / ss.capacity[st]
}
