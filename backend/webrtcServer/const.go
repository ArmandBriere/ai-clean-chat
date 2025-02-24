package webrtcserver

const (
	// WebRTC
	INPUT_SAMPLE_RATE = 48000
	MODEL_SAMPLE_RATE = 16000

	// Profanity
	PROFANITY_ANALYSIS_BUFFER_SIZE = 7
)

const LLM_PROMPT = `
You are a helpful assistant aiding a user in understanding transcribed audio.
The provided text contains profanity.

In under 20 words, explain why the text is profane.
Focus on the reason behind the profanity without explicitly stating "this is profane" or similar phrases.

Be concise, direct, and educational to help the user learn.
`
