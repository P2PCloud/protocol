package broker

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/p2pcloud/protocol"
)

func (b *Broker) BookVM(offerIndex, seconds int) error {
	tx, err := b.session.BookVM(big.NewInt(int64(offerIndex)), big.NewInt(int64(seconds)))
	if err != nil {
		return err
	}

	return b.waitForTx(tx.Hash())
}

func (b *Broker) GetUsersBookings() ([]protocol.VMBooking, error) {
	if err := b.setDecimals(); err != nil {
		return nil, err
	}

	bookings, err := b.session.FindBookingsByUser(b.transactOpts.From)
	if err != nil {
		return nil, err
	}
	var result []protocol.VMBooking
	for _, booking := range bookings {
		result = append(result, protocol.VMBooking{
			VmTypeId:   int(booking.VmTypeId.Int64()),
			PPS:        b.amountToCoins(booking.PricePerSecond),
			Miner:      &booking.Miner,
			Index:      int(booking.Index.Int64()),
			User:       &booking.User,
			BookedAt:   int(booking.BookedAt.Int64()),
			BookedTill: int(booking.BookedTill.Int64()),
		})
	}
	return result, nil
}

func (b *Broker) GetBooking(index int) (*protocol.VMBooking, error) {
	if err := b.setDecimals(); err != nil {
		return nil, err
	}

	booking, err := b.session.GetBooking(uint64(index))
	if err != nil {
		return nil, err
	}

	if booking.Miner.Hex() == "0x0000000000000000000000000000000000000000" {
		return nil, fmt.Errorf("booking %d not found", index)
	}

	return &protocol.VMBooking{
		VmTypeId:   int(booking.VmTypeId.Int64()),
		PPS:        b.amountToCoins(booking.PricePerSecond),
		Miner:      &booking.Miner,
		Index:      int(booking.Index.Int64()),
		User:       &booking.User,
		BookedAt:   int(booking.BookedAt.Int64()),
		BookedTill: int(booking.BookedTill.Int64()),
	}, nil
}

func (b *Broker) GetTime() (int, error) {
	t, err := b.session.GetTime()
	if err != nil {
		return 0, err
	}
	return int(t.Int64()), nil
}

func (b *Broker) GetMinersBookings() ([]protocol.VMBooking, error) {
	if err := b.setDecimals(); err != nil {
		return nil, err
	}

	bookings, err := b.session.FindBookingsByMiner(b.transactOpts.From)
	if err != nil {
		return nil, fmt.Errorf("error executing FindBookingsByMiner: %v", err)
	}
	var result []protocol.VMBooking
	for _, booking := range bookings {

		result = append(result, protocol.VMBooking{
			VmTypeId:   int(booking.VmTypeId.Int64()),
			PPS:        b.amountToCoins(booking.PricePerSecond),
			Miner:      &booking.Miner,
			Index:      int(booking.Index.Int64()),
			User:       &booking.User,
			BookedAt:   int(booking.BookedAt.Int64()),
			BookedTill: int(booking.BookedTill.Int64()),
		})
	}
	return result, nil
}

func (b *Broker) AbortBooking(index uint64, abortType protocol.AbortType) error {
	tx, err := b.session.AbortBooking(index, abortType.ToSolidityType())
	if err != nil {
		return err
	}

	return b.waitForTx(tx.Hash())
}

func (b *Broker) ClaimExpired(index uint64) error {
	tx, err := b.session.ClaimExpired(index)
	if err != nil {
		return err
	}

	return b.waitForTx(tx.Hash())
}

func (b *Broker) ExtendBooking(index uint64, secs int) error {
	tx, err := b.session.ExtendBooking(index, big.NewInt(int64(secs)))
	if err != nil {
		return err
	}

	return b.waitForTx(tx.Hash())
}

func (b *Broker) GetUserBookings() ([]protocol.VMBooking, error) {
	bb, err := b.session.GetUsersBookings(crypto.PubkeyToAddress(b.GetPrivateKey().PublicKey))
	if err != nil {
		return nil, err
	}

	result := make([]protocol.VMBooking, 0, len(bb))

	for i := range bb {
		result = append(result, protocol.VMBooking{
			VmTypeId:   int(bb[i].VmTypeId.Int64()),
			PPS:        b.amountToCoins(bb[i].PricePerSecond),
			Miner:      &bb[i].Miner,
			Index:      int(bb[i].Index.Int64()),
			User:       &bb[i].User,
			BookedAt:   int(bb[i].BookedAt.Int64()),
			BookedTill: int(bb[i].BookedTill.Int64()),
		})
	}

	return result, nil
}

// func (b *Broker) GetMinerUrl(address *common.Address) (string, error) {
// 	urlBytes, err := b.session.GetMinerUrl(*address)
// 	if err != nil {
// 		return "", err
// 	}
// 	url, err := converters.Bytes32ToUrl(urlBytes)
// 	return url, err
// }

// func (b *Broker) SetMinerUrlIfNeeded(newUrl string) error {
// 	oldUrl, err := b.GetMinerUrl(&b.transactOpts.From)
// 	if err != nil {
// 		return err
// 	}
// 	if oldUrl == newUrl {
// 		return nil

// 	}

// 	urlBytes, err := converters.UrlToBytes32(newUrl)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = b.session.SetMunerUrl(urlBytes)
// 	return err
// }
