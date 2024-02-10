TODO:

- Write checker.
- Reinstate WriteOutput (and perhaps omit log to console, or maybe not; it's convenient and shows the end of what might be a very long file; at the audit, I should explain that the final printout should not count as part of the program time; include a note in the README to that effect).

- Refactor: trim process struct, neaten buildOutput, schedule, and graph.
- Scrutinize logic.

# stock-exchange-sim

0. [Brief](#-brief)
1. [Setup](#1-setup)
2. [Usage](#2-usage)
3. [Audit](#3-audit)
4. [Research](#3-research)
5. [Strategy](#4-strategy)

## 0. Brief

According to the [instructions](https://github.com/01-edu/public/tree/master/subjects/stock-exchange-sim), we need to write a program, `stock`, that takes two command-line arguments.

This first is the path to a text file, which they call a "configuration file", formatted in a certain way. Our program should derive information, from this file, about a project. The file will contain a list of resources and quantities of each available at the start of the project; a list of processes with information about which resources they consume, in what quantity, and what products they produce and in what quantity; also, how long each process takes. Finally, it will specify the goal of the project, in the form of an item whose production should be "optimized", i.e. maximized.

The configuration file can also specify that time should be optimized too. In the case of non-renewable resources, we take this to mean that time should be minimized given that the maximum amount of the goal should be produced.

Comments in some of the examples suggest the possibility of multiple resources to optimize, but none of the examples actually realize that possibility. The line of the examples that cites the goal always has the format `optimize:(<stock_name>)` or `optimize:(time;<stock_name>)`. Although the instructions speak of "elements" to optimize, the format they specify is `optimize:(<stock_name>|time)`. No indication is given of how one would decide between conflicting goals. We could show precedence by the order they're listed in, but, for now, have taken the easier path of assuming only one stock item is to be maximized.

The second argument is an integer representing the maximum number of seconds the program must take to execute.

There are two types of task: those that can continue indefinitely thanks to renewable resources, and those that have a finite objective. (Given that the number of instances of a task that can be scheduled at once is only limited by resources and precedence relations, time optimization would seem to only trivially relevant when resources are not renewable.)

Given a configuration file `build.txt`, our program, `stock`, should produce a text file `build.log`, consisting of a schedule: a list of processes (possibly including several instances of the same process, possibly overlapping), the statement "No more process duable at", followed by an integer one unit greater than the duration of whole project, and a list of stock (resources and products) left at the end.

## 1. Setup

Optional: to build an executable file of the `stock` program, navigate into the `stock-exchange-sim` folder and run the command `go build -o stock`.

## 2. Usage

Run `go . simple.txt` (or, after building an executable, `./stock examples/simple`) to create a schedule for the example called `simple`.

Run `go . simple 10` to specify that the program should not take longer than 10 seconds.

## 3. Audit

You'll find the configuration files and logs for the given examples in the `examples` folder, together with the two examples we were required to create: `zen` (finite) and `matryushka` (infinite).

Exact outputs may vary from those suggested in the audit questions, especially where time is not to be optimized, since then there is less constraint on how soon they can be scheduled. Thus, for `seller`, the audit suggestion takes a more leisurely approach, whereas our program schedules processes as soon as the precedence relations allow, because why not?

Please note that we've chosen to interpret the time parameter as marking when to end the schedule function. Writing the output file and printing the result to the terminal happen after that. This seemed like the most natural interpretation of the instructions, particularly as printing to the terminal is an optional extra for ease of viewing long outputs. It has no bearing on the result of the audit, which just asks you to confirm that fewer processes are performed in 0.003s than 1s for the example `fertilizer`.

## 4. Research

As recommended, we consulted [PM Knowledge Center](https://www.pmknowledgecenter.com), a collection of resources on "Project Management and Dynamic Scheduling". We found further background reading necessary to fill in the gaps in the explanations there: in particular, [Kolisch (1994)](https://www.econstor.eu/bitstream/10419/155418/1/manuskript_344.pdf). These sources describe what's known as a Resource Constrained Project Scheduling Problem. The heuristic type of solution our instructions direct us towards is called Priority Rule Based Scheduling.

Before going into detail, we should note that the above sources (or Kolish, at least, who gives more detailed algorithms) assume that each task can only be performed once per project, whereas our program is expected to deal with cases where tasks can and should be performed more than once (in succession or simultaneously), if resorces allow, to optimize what needs optimizing. This meant that we couldn't directly apply either of the proposed scheduling methods.

In Priority Rule Based Scheduling, a graph of precedence relationships is drawn up: that is, a graph where tasks are nodes, and a directed edge from A to B means that commencement of B depends directly on completion of A. Activities are numbered in such a way that successors always have a greater activity number than their predecessors. A priority rule is chosen. Then a schedule is generated according to one of two schemes:

- Serial

- Parallel

A SERIAL schedule generation scheme with N tasks takes N steps. One task is chosen, at each step, from the set of available tasks and moved to the set of completed tasks. (A task is available if it is the direct the successor to a completed task, and current resources suffice to perform it.) If multiple tasks are available, one is chosen according to the priority rule. If several have equal priority, the one with the lowest activity number is selected.

A PARALLEL schedule generation scheme with N tasks takes at most N steps. At each step, we schedule zero or more activities. Tasks are partitioned into completed, in progress, and available. The schedule time associated with a step is calculated as the earliest completion time of the tasks that were in progress during the previous step. Tasks whose finish time is equal to the schedule time are moved from the set of tasks in progress to the set of completed tasks. This may make other tasks available. As long as tasks are available, they're chosen one by one, in order as in a serial scheme, and started at the current schedule time, then we move on to the next step. The algorithm terminates when all tasks are completed or in progress.

## 5. Strategy

After all that, neither scheme quite works for us, given the different underlying assumptions of our project: multiple instances of a task schedulable, possibly simultaneously. But we can take inspiration from them.

We start with the simplifying assumption that, as in our examples, tasks can have multiple predecessors but only one sucessor. We also make some "good faith" assumptions about the configuration file, such as the absence of processes that don't contribute towards the goal. A more robust scheme would need to deal with such cases.

For a finite project, let the `count` of a task be the number of times it's scheduled to be performed. For a project that is infinite, due to renewable resourses, the `count` of a task will be the number of times it's scheduled to be performed per iteration of a loop that runs till interupted by a signal from the timer.

Set all `count`s to zero initially. Also give each process a field `minCount` that will be used to initialize the count. `minCount` will be of a home-made type, rational, representing a rational number. Set the `minCount` of each task to 1 initially. Then procede backwards from final tasks (defined as those that directly produce the target item) to initial ones (those that can be performed immediately). At each iteration, identify predecessors and successors, and set the `minCount` of a process equal to the `minCount` of its successor times the quantity that its successor needs of the item by which they're linked, divided by the quantity that it produces of that item.

Now, define `maxCount` as the least common multiple of the denominators of all the `minCount`s, and set the `count` of a process equal to the numerator of `maxCount` times its `minCount`.

At each stage of scheduling, we check whether a process can be performed `count` number of times, given the resources. If not, it can't be scheduled. If so, schedule it as soon as possible, given the durations and start times of its precursors. For finite projects, the end time is the start time of the final process plus its duration. For cycling projects, the provisional end time is continually updated as tasks are scheduled, and returned along with the schedule when the timer signals to finish.
