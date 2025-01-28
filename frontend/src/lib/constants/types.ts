export type StreamingOfferMessage = {
  type: string;
  payload: RTCSessionDescriptionInit;
};

export type StreamingAnswerMessage = {
  type: string;
  payload: RTCSessionDescriptionInit;
};

export type StreamingIceCandidateMessage = {
  type: string;
  payload: RTCIceCandidate;
};
