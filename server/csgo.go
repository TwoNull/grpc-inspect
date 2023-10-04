package main

import (
	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	pb "github.com/twonull/grpc-inspect/proto"
)

const AppId = 730

type CSGO struct {
	client *steam.Client
}

type (
	ClientReady struct{}
)

func InitClient(client *steam.Client) *CSGO {
	t := &CSGO{client}
	client.GC.RegisterPacketHandler(t)
	return t
}

func (t *CSGO) SetPlaying(playing bool) {
	if playing {
		t.client.GC.SetGamesPlayed(AppId)
	} else {
		t.client.GC.SetGamesPlayed()
	}
}

func (t *CSGO) SendHello() {
	t.client.GC.Write(gamecoordinator.NewGCMsgProtobuf(AppId, uint32(protobuf.EGCBaseClientMsg_k_EMsgGCClientHello), &protobuf.CMsgClientHello{}))
}

func (t *CSGO) InspectItem(s, a, d, m uint64) {
	InspectPayload := &protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockRequest{
		ParamS: &s,
		ParamA: &a,
		ParamD: &d,
		ParamM: &m,
	}

	t.client.GC.Write(gamecoordinator.NewGCMsgProtobuf(AppId, uint32(protobuf.ECsgoGCMsg_k_EMsgGCCStrike15_v2_Client2GCEconPreviewDataBlockRequest), InspectPayload))
}

func (t *CSGO) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	switch packet.MsgType {
	case uint32(protobuf.EGCBaseClientMsg_k_EMsgGCClientConnectionStatus):
		t.handleConnectionStatus(packet)
	case uint32(protobuf.ECsgoGCMsg_k_EMsgGCCStrike15_v2_Client2GCEconPreviewDataBlockResponse):
		t.handleEconPreviewDataBlockResponse(packet)
	}
}

func (t *CSGO) handleConnectionStatus(packet *gamecoordinator.GCPacket) {
	t.SendHello()
	t.SetPlaying(true)
}

func (t *CSGO) handleEconPreviewDataBlockResponse(packet *gamecoordinator.GCPacket) {
	data := &protobuf.CMsgGCCStrike15V2_Client2GCEconPreviewDataBlockResponse{}
	packet.ReadProtoMsg(data)
	dataInfo := data.GetIteminfo()

	itemInfo := &pb.ItemInfo{
		AccountId:          dataInfo.GetAccountid(),
		ItemId:             dataInfo.GetItemid(),
		DefIndex:           dataInfo.GetDefindex(),
		PaintIndex:         dataInfo.GetPaintindex(),
		Rarity:             dataInfo.GetRarity(),
		Quality:            dataInfo.GetQuality(),
		PaintWear:          dataInfo.GetPaintwear(),
		PaintSeed:          dataInfo.GetPaintseed(),
		KilleaterScoreType: dataInfo.GetKilleaterscoretype(),
		KilleaterValue:     dataInfo.GetKilleatervalue(),
		CustomName:         dataInfo.GetCustomname(),
		Inventory:          dataInfo.GetInventory(),
		Origin:             dataInfo.GetOrigin(),
		QuestId:            dataInfo.GetQuestid(),
		DropReason:         dataInfo.GetDropreason(),
		MusicIndex:         dataInfo.GetMusicindex(),
	}
	stickers := dataInfo.GetStickers()
	for _, sticker := range stickers {
		itemInfo.Stickers = append(itemInfo.Stickers, &pb.Sticker{
			Slot:      sticker.GetSlot(),
			StickerId: sticker.GetStickerId(),
			Wear:      sticker.GetWear(),
			Scale:     sticker.GetScale(),
			Rotation:  sticker.GetRotation(),
			TintId:    sticker.GetTintId(),
		})
	}
	t.client.Emit(itemInfo)
}
