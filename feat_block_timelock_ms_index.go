package iotago

import (
	"encoding/json"
	"fmt"

	"github.com/iotaledger/hive.go/serializer"
)

// TimelockMilestoneIndexFeatureBlock is a feature block which puts a time constraint on an output depending
// on the latest confirmed milestone index X:
//	- the output can only be consumed, if X is bigger than the one defined in the timelock.
type TimelockMilestoneIndexFeatureBlock struct {
	MilestoneIndex uint32
}

func (s *TimelockMilestoneIndexFeatureBlock) Deserialize(data []byte, deSeriMode serializer.DeSerializationMode) (int, error) {
	return serializer.NewDeserializer(data).
		AbortIf(func(err error) error {
			if deSeriMode.HasMode(serializer.DeSeriModePerformValidation) {
				if err := serializer.CheckTypeByte(data, FeatureBlockTimelockMilestoneIndex); err != nil {
					return fmt.Errorf("unable to deserialize timelock milestone index feature block: %w", err)
				}
			}
			return nil
		}).
		Skip(serializer.SmallTypeDenotationByteSize, func(err error) error {
			return fmt.Errorf("unable to skip timelock milestone index feature block type during deserialization: %w", err)
		}).
		ReadNum(&s.MilestoneIndex, func(err error) error {
			return fmt.Errorf("unable to deserialize milestone index for timelock milestone index feature block: %w", err)
		}).
		Done()
}

func (s *TimelockMilestoneIndexFeatureBlock) Serialize(_ serializer.DeSerializationMode) ([]byte, error) {
	return serializer.NewSerializer().
		WriteNum(s.MilestoneIndex, func(err error) error {
			return fmt.Errorf("unable to serialize timelock milestone index feature block milestone index: %w", err)
		}).
		Serialize()
}

func (s *TimelockMilestoneIndexFeatureBlock) MarshalJSON() ([]byte, error) {
	jTimelockMilestoneFeatBlock := &jsonTimelockMilestoneIndexFeatureBlock{MilestoneIndex: int(s.MilestoneIndex)}
	jTimelockMilestoneFeatBlock.Type = int(FeatureBlockTimelockMilestoneIndex)
	return json.Marshal(jTimelockMilestoneFeatBlock)
}

func (s *TimelockMilestoneIndexFeatureBlock) UnmarshalJSON(bytes []byte) error {
	jTimelockMilestoneFeatBlock := &jsonTimelockMilestoneIndexFeatureBlock{}
	if err := json.Unmarshal(bytes, jTimelockMilestoneFeatBlock); err != nil {
		return err
	}
	seri, err := jTimelockMilestoneFeatBlock.ToSerializable()
	if err != nil {
		return err
	}
	*s = *seri.(*TimelockMilestoneIndexFeatureBlock)
	return nil
}

// jsonTimelockMilestoneIndexFeatureBlock defines the json representation of a TimelockMilestoneIndexFeatureBlock.
type jsonTimelockMilestoneIndexFeatureBlock struct {
	Type           int `json:"type"`
	MilestoneIndex int `json:"milestoneIndex"`
}

func (j *jsonTimelockMilestoneIndexFeatureBlock) ToSerializable() (serializer.Serializable, error) {
	return &TimelockMilestoneIndexFeatureBlock{MilestoneIndex: uint32(j.MilestoneIndex)}, nil
}