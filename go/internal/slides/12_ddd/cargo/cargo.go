package cargo

import "errors"

type trackingID struct {
	value string
}

var ErrEmptyTrackingID = errors.New("empty tracking id")

func NewTrackingID(value string) (trackingID, error) {
	if value == "" {
		return trackingID{}, ErrEmptyTrackingID
	}
	return trackingID{
		value: value,
	}, nil
}

func (id trackingID) String() string {
	return id.value
}

func (id trackingID) Value() string {
	return id.value
}

type Cargo struct {
	TrackingID trackingID
}
