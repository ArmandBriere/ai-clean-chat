interface IceCandidateMessage {
    type: 'iceCandidate';
    candidate: RTCIceCandidateInit;
}

interface OfferMessage {
    type: 'offer';
    sdp: string;
}

interface AnswerMessage {
    type: 'answer';
    sdp: string;
}


interface StreamingMessage {
    type: 'streaming';
    isStreaming: boolean;
}

type WebSocketMessage = IceCandidateMessage | OfferMessage | AnswerMessage;
