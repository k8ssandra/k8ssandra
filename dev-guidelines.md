# Developer Guidelines

## Values

As a developer community, it's worth reminding ourselves of our broader community [values](https://k8ssandra.io/community/values/) and considering how they should impact the way we interact with each other throughout the development process.

Here's a quick review of our community values:

1. Value time
2. Beginner's mind
3. Openness
4. Positivity
5. Plans

Keep these values in mind as we discuss the responsibilities throughout the development and review process.

## Pull request/code review guidelines

### Goals

A wealth of great content has been written over the years about why code review is important.  We won't try to list every aspect of the concept, but we'll focus on a few key outcomes that our community desires from the pull request process.

#### Improvement

At the end of the day, the pull request and code review process is about improvement.

* Improvement of the specific code being submitted

* Improvement of the overall codebase 

* Personal improvement of the author and reviewer, in both technical and non-technical ways

#### Knowledge sharing

A major outcome of the review process should be learning.  Learning from each other not only coding practices and approaches related to "how" something has been implemented, but also context and insight into "why" something has been implemented.

It becomes a historical record for our future selves, so it's not only about sharing what we know right now, but it's about reminding ourselves of our decisions and thought processes later.

#### Communication and community building

The review process one of the biggest touch points that project maintainers have both with each other and with other contributors.  This process is an opportunity each time to improve our communication skills and help invite others into the community.

### Roles and responsibilities

There are multiple roles involved in the process of taking code from development to acceptance into the project:

1. Team
2. Author
3. Reviewer

Each role has their own responsibilities to be upheld to make the process productive, efficient, and positive for all involved.

#### Team

**Foster Collaboration** *[Beginner's mind, Openness, Positivity]* 

In order for individuals to be successful in the process, the team and community as a whole must embrace creating an environment focused on nonjudgmental collaboration.  To do this, we expect everyone to:

* Keep a positive tone -- it's totally OK to use emojis :wink:

* Assume the best in others -- asynchronous communication is hard (that's much of the reason for these guidelines), assume first that no matter the potentially undesired tone used, others are trying to help not harm in this process

* Default to questions, not demands or accusations -- this is how we learn, together

* Have empathy -- you've been on the other side before, remember that, neither role is easy

* Compromise -- balance the ability to deliver with the desire for quality

* Breakout -- acknowledge when getting together to talk `face-to-face` is the best way to resolve discussions

**Account For Reviews** *[Value time, Plans]*

Reviews are not free and should be planned for as we think about both our short and long term commitments both individually and as a project.

**Learn** *[Beginner's mind, Value time]*

Getting the most out of our time means learning from the observations and outcomes of the process.  It's the responsibility of the full team to learn from reviews and improve both themselves as well as the process itself.

Remember that this process is iterative, we will repeat it many-many...many times in the lifetime of the project and each of those is a chance to do it in a better way.  When you find those better ways, come back here and update this guide :wink:.

**Speak Up** *[Openness]*

If you feel like something isn't right in this process, ask, clarify, speak up, help us fix it.

#### Author

**Early Engagement** *[Value time, Openness]*

This is a responsibility that really begins before the PR itself has ever been created.  Seek to find and involve reviewers early in the design/development process.  We'll discuss the process of identifying reviewers more later, just keep in mind that even quick conversations about requirements or tradeoffs can give a reviewer a great bit more context during review later.

**Curate Requests** *[Value time, Plans]*

Pull requests take planning, invest in that aspect early in the process.  This piece is so important that we'll discuss it again later, but keep a few key concepts in mind when building a request:

* Focused and smaller pull requests are easier to review, discuss, and get merged

* Take opportunities to separate large requests into staged evolutions where possible

* Stage commits to help guide reviewers through a request where possible

* Take advantage of the project's pull request template, it should serve as a guide to adding context to a request

* Leave comments within the review to mark particular areas of focus for reviewers and share other code-level specific context that might be useful before others begin review

**Self-Review** *[Value time]*

The first reviewer of a pull request should always be the author.  Make a habit of self-review.  It can help spot early issues and avoid needless back and forth with a reviewer.  Use a `draft` pull request to indicate that something isn't yet ready for review.

**Be Patient, Listen, Respond** *[Beginner's mind, Value time, Openness, Positivity]*

Be patient when asking for other's time to do reviews, but also know that they have a responsibility to be responsive as well.  Don't be afraid to ask for help and feedback, but be understanding of other priorities they have as well.

When you get feedback, listen to it, but know that the reviewer's opinion is not the end of the conversation.  Remember that this process is an opportunity for the reviewer to learn as well, they will likely not have the same level of depth and context around a change that you the author have; guide them and teach them -- but also be open to their feedback.  Keep the beginner's mind and positivity front and center.

#### Reviewer

**Commit Time** *[Value time, Plans]*

Contribution is the life blood of the project, both from other maintainers and from external contributors, don't forget this.  All of our time is valuable, but as a steward of the project, we ask that you plan to commit time to review.  Take the time in that review to be thoughtful and meaningful with your feedback.  Take the time to shape that feedback in a way that fits with the overall values of the project.

This can extend to time needed to pull and test the code being contributed and it's understood at a project level that this can be significant and you are empowered to take that time.

**Be Patient, Listen, Respond** *[Beginner's mind, Value time, Openness, Positivity]*

Much like the same reponsiblity shared by the author, you the reviewer are expected to be patient with an author as they too juggle priorities to respond to your questions and feedback and then hear out their point of view when they do.  Support a two way conversation and keep in mind that you too can learn from this process.

**Decline The Role (If Necessary)** *[Openness, Value time]*

If you do not believe that you are the right person to review a particular request, you are empowered to politley decline and offer alteratives.  This might happen because you don't currently have time or perhaps you don't believe you have the background and knowledge necessary to make productive decisions about a proposed set of changes.

**Favor Approval Over Perfection** *[Value time, Positivity]*

Having a bias towards forward progress doesn't have to sacrifice quality.  It does mean acknowledging that code is rarely perfect and does not have to be in order to push the project in a positive direction.

***Reviewers should favor approval when a pull request has reached a state where it is a positive improvement on the overall codebase, even if there is room for further improvement.***  Reviewers should keep in mind that a contribution doesn't need to be done in the exact form and fashion that they would have done it to positively impact our project.

**Reject Early (If Necessary)** *[Value time, Openness]*

There will be times where as the reviewer you have to make a difficult decision that a request should be rejected.  This may happen for many reasons, when it does, try to do it early in the process and be clear about the reasons for that.  Actively work to avoid situations where many rounds of changes have been made and yet the mark has still not be hit.  When that happens, the mark maybe not have been properly set or communicated.  Also seek to be thankful for the time an author has committed to the work, even if it ends in an unsuccessful state, it can still be a positive experience for that contributor to learn from if handled well.

**Ask Questions** *[Beginner's mind, Openness, Positivity]*

Default to asking questions.  Ask "how" something works or "why" something was done if it's not clear from the code and context given with the PR.  Instead of stating how you would have done something, ask if an alternative approach was considered.  Don't assume and don't imply, engage in a conversation.

**Constructive Feedback** *[Openness, Positivity]*

This is perhaps the biggest thing you, as the reviewer, are responsible for.  In general, it encompasses the expectation of behavior across the other responsibilities as well.

In all of our interactions as reviewers we should strive to be both respectful and constructive.  Keep some of these things in mind when reviewing:

* Resist feedback stated as a demand or accusation -- instead, seek to start discussion, ask questions

* Be direct and explain the "why" -- it's not enough to simply say a change isn't acceptable, explain your thought process and concerns, guide the author to the place that a change set needs to be for acceptance

* Give examples -- provide examples, documentation, or in-review suggestions about how you think something could be improved, teach don't criticize

**Stay Focused** *[Value time, Plan]*

Focus on the changes proposed and resist the urge to comment on other pieces of the surrounding code.  That feedback is also important to the project, but its place is likely to be found in the creation of a new issue.

## Finding reviewers

Every change within our projects require the review of at least one maintainer - so it's critical that reviewers can be identified and committed too early in the process.  These reviewers should be seen less as guardians at the gates and more as guides for the journey.  They can become partners in design and even implementation along the way, sponsors for your goals as an author of a change.

Identifying a reviewer is a shared responsibility of the author and the team.  The author should identify potential members of the community that are a good match for the type of change being planned and when that knowledge isn't immediately available to the author personally, the team should help find those people.  Reach out and engage people -- ask for help, it's a good thing. 

If you're an author in need of such help, reach out to the team on [discord](https://discord.gg/af82SnxzTm) in the #k8ssandra-dev channel.

### How many reviewers?

But, how many reviewers are needed on a pull request?  ***Generally speaking, the answer is one (1)***.  It's important for the maintainers of the project to trust each other to make good decisions for the project and also to ask for help when they feel that they are not equipped in a particular scenario to do so.

The specific answer is certainly always more complicated than that, but we should try to avoid situations where reviewers "pile on" and inject complication late in the process, often without context.  

We should, as much as possible, respect a chain of review.  This is a process where initial reviewers are empowered to approve a pull request or to seek additional review from others, others who might have more experience with a particular aspect of the code.  The subsequent reviews do not always have to be comprehensive in nature, but can be supplementary to the initial review.  Again, the initial reviewer is empowered to approve.  If the initial reviewer doesn't believe that they can do that they should ask for help and potentially even shadow another reviewer simply to gain some insight into a new area of code and functionality.

## Pull request creation

There is an art to creating a successful pull request.  We will not attempt to cover every element of that art here, but will highlight a few key aspects and ask you to leverage the guidance provided in our pull request template to help along the way.

* Focused -- ideally each pull request should be small in scope and focused on solving a single issue, this is a key aspect of keeping the review process manageable and efficient

* Descriptive title -- pull requests should be named in a way that is understandable and later searchable

* Associated with an issue -- pull requests should be associated with an issue that it is addressing

* Context -- a pull request should include sufficient context for a reviewer to understand:
  * What problem is being solved
  * Why the problem needs to be solved
  * How to reproduce the problem
  * How the problem is being solved
  * How to test and verify that the problem has been solved

* Include tests -- frequently the cause of a rejected or delayed PR is insufficient testing being included

## Credit where credit is due

We are certainly not the first community to discuss what makes a good pull request and how to build a successful and positive environment for code reviews to happen.  Here are some of the resources that we found impactful as we went through this process ourselves:

* [Writing A Good Pull Request](https://developers.google.com/blockly/guides/modify/contribute/write_a_good_pr)
* [The Standard of Code Review](https://google.github.io/eng-practices/review/reviewer/standard.html)
* [Bringing A Healthy Code Review Mindset To Your Team](https://www.smashingmagazine.com/2019/06/bringing-healthy-code-review-mindset/)
* [Pull Request Etiquette for Reviewers and Authors](https://betterprogramming.pub/pull-request-etiquettes-for-reviewer-and-author-f4e80360f92c)
* [Are Your Pull Requests Hard To Review? 5 Tips To Make Them Easier](https://betterprogramming.pub/are-your-pull-requests-hard-to-review-5-tips-to-make-them-easier-b6759b0749e8)




