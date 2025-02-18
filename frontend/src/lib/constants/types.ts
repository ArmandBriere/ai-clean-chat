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

export type AnalyzedMessage = {
  uuid: string;
  text: string;
  profanityScore: number;
};

export type LLMAnalysis = {
  analysis: string;
  userMessage: string;
};
