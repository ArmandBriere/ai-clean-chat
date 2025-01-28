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
