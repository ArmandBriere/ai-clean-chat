# Design

The design of the application was developed by [Osedea](https://osedea.com) and is heavily inspired by the design of [Google Meet](https://meet.google.com/).

The objective of the design is to augment the user experience by providing a live transcript with a summary of the conversation and the censored words.

## User Interface

The user interface is composed of the following elements:

- A landing page that allows the user to create a new room or join an existing one
- A video page that displays the video chat and the transcript
- A end page that displays the summary of the conversation and the censored words

The design is for full screen first.

### Home page

![Home page](./assets/home.png)

### Video page

Video page without transcription:

![Video page 01](./assets/video_01.png)

Admin can enable the transcription:

![Video page 02](./assets/video_02.png)

Admin can enable profanity filter:

![Video page 03](./assets/video_03.png)

Transcription displayed:

![Video page 04](./assets/video_04.png)

Profanity detected:

![Video page 05](./assets/video_05.png)

![Video page 06](./assets/video_06.png)

List of profanity:

![Video page 07](./assets/video_07.png)

### Summary page

End page with summary of the profanities:

![End page](./assets/summary.png)

End modal with explanation and score:

![End modal](./assets/end-modal.png)
