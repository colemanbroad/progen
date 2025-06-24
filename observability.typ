#set quote(block:true)
#show link: underline 


= OTEL soup

The #link("https://opentelemetry.io/docs/concepts/signals/")[OpenTelemetry] standard defines a bunch of terms related to how we observe state in distributed systems:

#quote[Another thing you’ll note is that each Span looks like a structured log. That’s because it kind of is! One way to think of Traces is that they’re a collection of structured logs with context, correlation, hierarchy, and more baked in. However, these “structured logs” can come from different processes, services, VMs, data centers, and so on. This is what allows tracing to represent an end-to-end view of any system.]


- *structured log* (aka slog, wide event) : has a name, time, other key-val pairs.

But OTEL defines the following terms, which add requirements to the slog:

- *signals* : container term for any of the below

- *traces* : The path of a request through your application.

- *metrics* : A measurement captured at runtime.
The slog for a metric is a *metric event*. 
A *meter* creates *metric instruments*, capturing measurements about a service at runtime. Meters are create with *Meter providers*. *Metric exporters* send metric data to a consumer.

#quote[In OpenTelemetry measurements are captured by metric instruments. A metric instrument has
- Name
- Kind
- Unit (optional)
- Description (optional)
]

- *logs* : 
- *baggage* : 
- *events* : 
- *profiles* : 

What are *spans*? Referenced on HN. 
What is *telemetry*?

[traces] Traces have a `trace_id` which is shared across all slogs created in response to the same request.
And they have a `span_id` and `parent_id` (also a span_id) which allow strcturing slogs into a tree.
A root span has `parent_id = null`. 

= References

- https://opentelemetry.io/docs/concepts/signals/metrics/
- https://coroot.com/blog/opentelemetry-for-go-measuring-the-overhead/ 
- https://news.ycombinator.com/item?id=20847352 Rust open source tool for metrics, logs, traces.
- https://news.ycombinator.com/item?id=32733620 Open source tool for dist metrics, logs, traces. 
- https://news.ycombinator.com/item?id=20375190 Logs vs Metrics. A false dichotomy.
- https://www.metaplane.dev/blog/the-origins-purpose-and-practice-of-data-observability#the-origins-of-data-observability
- https://www.metaplane.dev/blog/the-four-pillars-of-data-observability
- https://www.datadoghq.com/knowledge-center/observability/
- https://github.com/samber/slog-parquet
- https://grafana.com/
- https://thenewstack.io/honeycombs-charity-majors-go-ahead-test-in-production/
- https://www.honeycomb.io/blog/its-the-end-of-observability-as-we-know-it-and-i-feel-fine
- honeycomb https://www.honeycomb.io/blog/honeycomb-fit-software-development-lifecycle
- https://signoz.io/blog/jaeger-vs-prometheus/
- https://www.honeycomb.io/blog/dashboards-or-launchpads
- the idea of "shift left". 
- https://o11y.eu/ is a company that advises on observability.


#let l = "https://www.metaplane.dev/blog/data-observability-vs-software-observability"
== #link(l)[Data observability vs Software observability]

#quote[For a concept that no one can quite define, observability is hot.]

lol

References 

- Datadog,
- AppDynamics,
- New Relic,
- Grafana,
- Splunk,
- Sumo Logic,

#let l = "https://news.ycombinator.com/item?id=39529775"
== #link(l)[All you need are wide events] (i.e. structured logs).

https://isburmistrov.substack.com/p/all-you-need-is-wide-events-not-metrics

Judging from the HN comments it seems like *metrics* have come to mean "something that isn't subsampled" (i.e. it's value is recorded every time the line is hit).
Also associated with "metrics" is the idea that they are fundamentally near-continuous, but (sub)sampled by a heartbeat WE control,
e.g. a fixed 1Hz timer, e.g. "CPU util" is a metric.
They are measurements that can be performed at any point in the code, unlike e.g. sampling the value of a variable, which only makes sense during that var's lifetime. 
Comments include "Their architectures allow you to scale at low costs to PB or even exabyte scale monitoring". The idea of exabyte scale logging just seems wrong.

#let l = "https://x.com/mipsytipsy/status/1494856337632608261"
== #link(l)[Datadog screed by charity majors], in favor of OTel

Features this #link("https://x.com/apostolis09/status/1495042077285130245")[very good point in reply]. Not substantive, just Datadog are sneaky bastards subverting OTel.

#let l = "https://www.youtube.com/watch?v=ag2ykPO805M&ab_channel=GOTOConferences"
== #link(l)[Charity's talk at GOTO about observability 2.0].

So apparently all the complex OTEL terminology was a huge mistake (duh) and observability 2.0 makes it all radically simpler.
There's just structured logs and that's it!  
She recommends #link("https://stripe.com/blog/canonical-log-lines")[canonical log lines] which are one slog per request,
at the end of the request, with alllll the info (essentially the result of joining on the request id).

== https://stripe.com/blog/canonical-log-lines

Also https://news.ycombinator.com/item?id=20568634

The idea that you have a single, wide slog per request (or tied to some other important upsteam event's lifecycle). 

== https://www.hyperdx.io/blog/logs-matter-more-than-metrics

Accepts dichotomy where metrics give you plots and logs give you _nothing_, should be wide by default and used to better understand cause.
Thesis is that #quote[When an incorrect behavior emerges, logs are more likely to explain what happened than any metrics.]. 

= Observability as a version of medicine.

Imagine if human health worked this way.
You walk around all the time with a thermometer in your mouth and on your forehead, blood pressure on your arm and a glucose monitor on your finger. You've got a covid test up your nose and a needle that constantly draws blood to check for diabetes and immune issues.
These complex instruments must be maintained constantly to be in a working state!

NO!

You notice that you're sick because of symptoms that affect you.

That *don't* necessarily show up when you're at work. In fact they mostly show up when you're doing some other, more strenuous activity: picking up your child or a suitcase reveals problems in your knees and back. Reading at night reveals problems with
your vision. It's the weird, diverse activities outside of the routine of work that reveal problems. Sometimes you get lucky
and a problem is revealsed when you go for your annual physical! But mostly that's not what happens.
And how often do you hear about a physical that reveals a problem that never manifested any symptoms, but solving the problem involved a surgery that left a scar and required months of recovery? Or how often the drugs you take cancel each other out?

Experiment in the wild on real systems.
Test in prod.
That's the best bang-for-your-buck.

Is littering your code with assertions the same thing as logging?
No. For one, there's no storage cost every time an assertion is hit.
Only when they fail.

OK, how about this.
Instead of compiling release and debug versions of your code you compile BOTH and toggle back and forth dynamically.
Every `n` requests can pass through the debug version of your code, which logs everything.

= Test in Prod, Antithesis, or a new Third Way?

Antithesis' "global bit-level equality" version of determinism prevents us from re-running an input but with e.g. verbose logging toggled ON. This means we have to do clever, very expensive things like Cameron's idea to periodically inspect the state of the system on a new, short branch forming the original branch with periodic, dangling observe-only branches.
But a softer notion of equality might lead to Heisenbugs when we toggle logging ON and the bug happens to depend on a
e.g. the size-on-disk of folders in the parent dir, or perhaps a microsecond long timing window. 

Dave's thesis is that complex systems are fundamentally chaotic, and that seemingly-irrelevant changes to e.g. log level will
very often have a butterfly effect on something upstream of the bug.  

Here's the possible choices:

1. log everything and sort through the garbage when you crash to find the clues
2. mock everything and run in sim, with determinism & fault injection. when bugs appear travel back in time and *manually* investigate your system.
3. Log litte, mock nothing. Cache all responses from external services? cache results from expensive operations. When a bug is triggered re-run with the same inputs to see how likely the bug is to re-appear. Maybe increase the verbosity of logging? Maybe
isolate components? Maybe factor your code to make it easy to bisect problems? If turning on observability features makes the bug
go away (Heisenbug) then either (a) you are reading and acting based on your logs (b) you are reading/acting based on some other
piece of the system that's changed (e.g. hash of / size of the contents of system folders).

Betting that 3 is better than 2 is betting that the difficulty of mocking outweights the likelihood of running in to a heisenbug. 
Or that mocking services is more likely to hide a bug than the prob of having a real heisenbug.

--- --- --- --- --- --- --- --- ---

You could run with logging ON, but automatically delete logs if no properties are violated?
This means you can't debug unexpected failures.
And it doesn't fix the problem where you hit a property failure but didn't have sufficient logging!

