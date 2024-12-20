package webrtcserver

// PcmToFloat32 converts PCM samples to normalized float32 samples
func PcmToFloat32(pcm []int16) []float32 {
	floatSamples := make([]float32, len(pcm))
	for i, sample := range pcm {
		// Normalize from int16 (-32768 to 32767) to float32 (-1.0 to 1.0)
		floatSamples[i] = float32(sample) / 32768.0
	}
	return floatSamples
}
