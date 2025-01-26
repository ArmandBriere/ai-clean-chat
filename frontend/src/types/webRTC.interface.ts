const ICE_CANDIDATE = 'iceCandidate';
const OFFER = 'offer';
const ANSWER = 'answer';
const STREAMING = 'streaming';
const TRANSCRIPTION = 'transcription';

interface IceCandidateMessage {
  type: string;
  candidate: RTCIceCandidateInit;
}

interface OfferMessage {
  type: string;
  sdp: string;
}

interface AnswerMessage {
  type: string;
  sdp: string;
}

interface StreamingMessage {
  type: string;
  isStreaming: boolean;
}

type WebSocketMessage = IceCandidateMessage | OfferMessage | AnswerMessage;
