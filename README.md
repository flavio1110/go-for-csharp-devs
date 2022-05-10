# Go for C# developers

> NB! This is and always will be a work in progress!

Learning GO is not only about a new language, but also about a different philosophy from what we are used to.

This page contains some useful links and answers to questions I had a long the way. Hopefully, it can help you too.

> I'll keep adding more materials, details, and answers a long the way, but as anything in life, the learning process is endless, therefore, the material in here can become outdated. 

## Materials

### Official docs

- Effective Go - [https://go.dev/doc/effective_go](https://go.dev/doc/effective_go)
    - A document that gives tips for writing clear, idiomatic Go code. A must-read for any new Go programmer. It augments the tour and the language specification, both of which should be read first.
- Go Docs - [https://go.dev/doc/](https://go.dev/doc/)
    - Root for many useful documentations
- FAQ - [https://go.dev/doc/faq](https://go.dev/doc/faq)
    - Answers to common questions about Go.

### Courses

- The way to go - [https://www.educative.io/courses/the-way-to-go](https://www.educative.io/courses/the-way-to-go)
    - It's a course from educative.io, it goes from the basics concepts of the language to more advanced ones. It also brings very interesting insights into the differences between the approaches of Java/C# to Go.
    - It has many challenges in exercises you can do directly in the browser. I found it excellent because I didn't need to have any thing installed and I managed to do it in anything that has a browser.

### Videos

- Go in 100 seconds - [https://www.youtube.com/watch?v=446E-r0rXHI&t=38s](https://www.youtube.com/watch?v=446E-r0rXHI&t=38s)
    - Short introduction about Go
- Simplicity is complicated - [https://www.youtube.com/watch?v=rFejpH_tAHM](https://www.youtube.com/watch?v=rFejpH_tAHM)
    - Rob Pike talks about how Go is often described as a simple language. It is not, it just seems that way. Rob explains how Go's simplicity hides a great deal of complexity, and that both the simplicity and complexity are part of the design.
- Concurrency is not parallelism - [https://www.youtube.com/watch?v=qmg1CF3gZQ0&t=1582s](https://www.youtube.com/watch?v=qmg1CF3gZQ0&t=1582s)
    - Rob Pike talks about concurrency and how Go implements it
- Understanding channels - [https://www.youtube.com/watch?v=KBZlN0izeiY&t=1011s](https://www.youtube.com/watch?v=KBZlN0izeiY&t=1011s)
    - Channels provide a simple mechanism for goroutines to communicate, and a powerful construct to build sophisticated concurrency patterns. We will delve into the inner workings of channels and channel operations, including how they're supported by the runtime scheduler and memory
- Concurrency in Go - [https://www.youtube.com/watch?v=_uQgGS_VIXM&list=PLsc-VaxfZl4do3Etp_xQ0aQBoC-x5BIgJ](https://www.youtube.com/watch?v=_uQgGS_VIXM&list=PLsc-VaxfZl4do3Etp_xQ0aQBoC-x5BIgJ)
    - Playlist with few short videos about different components of concurrency in Go
    
## Experiments
The folder experiments will contain pieces of code I'm play around, and also examples related to the questions below.

## FAQ for C# developers

- For what should I use Go?
    - Go has several uses, since command line applications, web services (yes web apis too :) ), and container applications like Docker and Kubernetes.

- What text-editor/IDE should I use
    - Visual Studio Code has everything you need to develop your application, and it's free.
    It's one of the best text editors, the debugger experience is decent, and there are really good plugins to help out with producitivity. [The extension](https://code.visualstudio.com/docs/languages/go) developed by Google itself is a must have.

- Does Go compile to some sort of Intermediate Language?
    - Nope. There is no IL, the compiler generates an executable that can be executed [directly by the computer](https://getstream.io/blog/how-a-go-program-compiles-down-to-machine-code/).

- How can I write async code?
    - Yes! In go we don't deal with threads directly. Instead, we use something called goroutines. In some sense they are similar to a `Task` which is not a necessary a thread but represents some work to be done, which will be eventually scheduled, executed in a thread, and resumed.
    I'm not naive to try to explain it in a short answer, so check [this video](https://www.youtube.com/watch?v=qmg1CF3gZQ0&t=1582s) from Rob Pike about it.

- Do I need to use goroutines and channels to make IO calls?
- How can I implement an interface?
- What about Nullable objects?
- Is there anything similar to LINQ?
- What about dependency injection containers?
- Which framework should I use for writing web APIs?
- What is the preferred t way to write logs?
- How should I organize my project?
- What IDE should I use