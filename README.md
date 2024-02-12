# stock-exchange-sim

0. [Brief](#0-brief)
1. [Setup](#1-setup)
2. [Usage](#2-usage)
3. [Audit](#3-audit)
4. [Research](#4-research)
5. [Strategy](#5-strategy)
6. [Further](#6-further)

## 0. Brief

This is a project in the [01 Edu](https://01-edu.org/) system, introducing the idea of a [process chain](https://en.wikipedia.org/wiki/Event-driven_process_chain). It's an optional extra, at least for us at 01Founders in London, 2024. It can be done in any compiled language. We chose Go.

According to the [instructions](https://github.com/01-edu/public/tree/master/subjects/stock-exchange-sim), we need to write a program, `stock`, that takes two command-line arguments.

This first is the path to a text file, which they call a "configuration file", formatted in a certain way. Our program should derive information, from this file, about a project. The file will contain a list of resources and quantities of each available at the start of the project; a list of processes with information about which resources they consume, in what quantity, and what products they produce and in what quantity; also, how long each process takes. Finally, it will specify the goal of the project, in the form of an item whose production should be "optimized", i.e. maximized.

The configuration file can also specify that time should be optimized too. In the case of non-renewable resources, we take this to mean that time should be minimized provided that the maximum amount of the goal is produced.

Comments in some of the examples suggest the possibility of multiple resources to optimize, but none of the examples actually realize that possibility. The line of the examples that cites the goal always has the format `optimize:(<stock_name>)` or `optimize:(time;<stock_name>)`. Although the instructions speak of "elements" to optimize, the format they specify is `optimize:(<stock_name>|time)`. No indication is given of how one would decide between conflicting goals. We could show precedence by the order they're listed in, but, for now, have taken the easier path of assuming only one stock item is to be maximized.

The second argument is an integer representing the maximum number of seconds the program must take to execute.

There are two types of task: those that can continue indefinitely thanks to renewable resources, and those that have a finite objective.

Given a configuration file `build`, our program, `stock`, should produce a text file `build.log`, consisting of a schedule: a list of processes (possibly including several instances of the same process, possibly overlapping), the statement "No more process duable at", followed by an integer one unit greater than the duration of whole project, and a list of stock (resources and products) left at the end.

We should also make a checker that will check the processes listed in a log file and confirm that there are enough resources to perform each task listed at the specified start time.

## 1. Setup

Optional: to build an executable file of the `stock` program, navigate into the `stock-exchange-sim` folder and run the command `go build -o stock`.

## 2. Usage

Run `go . simple` (or, after building an executable, `./stock examples/simple`) to create a schedule for the example called `simple`.

Run `go . simple 10` to specify that the program should not take longer than 10 seconds.

We've chosen to implement the checker as part of the same program. To check `simple.log`, run `go run . -checker examples/simple examples/simple.log` (or `./stock -checker examples/simple examples/simple.log` as the case may be).

## 3. Audit

As mentioned, we implemented the checker as part of the main program. A boolean flag is used to select checker mode. See `main.go` for the code that deals with the flag and other arguments, and `checker.go` for the checker function itself.

You'll find the configuration files and logs for the given examples in the `examples` folder, together with the two examples we were required to create: `zen` (finite) and `matryushka` (infinite).

Exact outputs may vary from those suggested in the audit questions, especially where time is not to be optimized, since, in that case, there is less constraint on how soon tasks can be scheduled. Thus, for `seller`, the audit suggestion takes a more leisurely approach, whereas our program schedules processes as soon as the precedence relations allow, because why not?

Please note that we've chosen to interpret the time parameter as marking when to end the schedule function. Writing the output file and printing the result to the terminal happen after that. This seemed like the most natural interpretation of the instructions, particularly as printing to the terminal is an optional extra for ease of viewing long outputs. It has no bearing on the result of the audit, which just asks you to confirm that fewer processes are performed in 0.003s than 1s for the example `fertilizer`.

## 4. Research

Following the recommendation of the project description, we consulted [PM Knowledge Center](https://www.pmknowledgecenter.com), a collection of resources on "Project Management and Dynamic Scheduling". We found further background reading necessary to fill in the gaps in the explanations there: in particular, [Kolisch (1994)](https://www.econstor.eu/bitstream/10419/155418/1/manuskript_344.pdf). These sources describe what's known as a Resource Constrained Project Scheduling Problem. The heuristic type of solution our instructions direct us towards is called Priority Rule Based Scheduling.

Before going into detail, we should note that the above sources (or Kolish, at least, who gives more detailed algorithms) assume that each task can only be performed once per project, whereas our program is expected to deal with cases where tasks can and should be performed more than once (in succession or simultaneously), if resorces allow, to optimize what needs optimizing. This meant that we couldn't directly apply either of the proposed scheduling methods.

In Priority Rule Based Scheduling, a graph of precedence relationships is drawn up: that is, a graph where tasks are nodes, and a directed edge from A to B means that commencement of B depends directly on completion of A. Activities are numbered in such a way that successors always have a greater activity number than their predecessors. A priority rule is chosen. Then a schedule is generated according to one of two schemes:

- Serial

- Parallel

A SERIAL schedule generation scheme with N tasks takes N steps. One task is chosen, at each step, from the set of available tasks and moved to the set of completed tasks. (A task is available if it's the direct the successor to a completed task, and current resources suffice to perform it.) If multiple tasks are available, one is chosen according to the priority rule. If several have equal priority, the one with the lowest activity number is selected.

A PARALLEL schedule generation scheme with N tasks takes at most N steps. At each step, we schedule zero or more activities. Tasks are partitioned into completed, in progress, and available. The schedule time associated with a step is calculated as the earliest completion time of the tasks that were in progress during the previous step. Tasks whose finish time is equal to the schedule time are moved from the set of tasks in progress to the set of completed tasks. This may make other tasks available. As long as tasks are available, they're chosen one by one, in order as in a serial scheme, and started at the current schedule time, then we move on to the next step. The algorithm terminates when all tasks are completed or in progress.

## 5. Strategy

After all that, neither scheme quite works for us, given the different underlying assumptions of our project: multiple instances of a task schedulable, possibly simultaneously. But we can take inspiration from them.

We start with the simplifying assumption that, as in our examples, processes can have multiple predecessors but only one sucessor. We also make some "good faith" assumptions about the configuration file, such as the absence of processes that don't contribute towards the goal. A more robust scheme would need to deal with such cases.

Our program scheduled tasks by taking as many passes through the precedence graph as resources permit. Before the first pass, an array (in Go terms, a slice) of the tasks currently being considered is initialized with tasks that can start at once.

We define the `count` of a process as the number of times it needs to be scheduled to produce one unit of the target item. The `iterations` of a process will be the number of times it is eventually scheduled.

Set all `count`s to zero initially. Also give each process a field `minCount` that will be used to initialize the count. `minCount` will be of a home-made type, rational, representing a rational number. Set the `minCount` of each task to 1 initially. Then procede backwards from final tasks (defined as those that directly produce the target item) to initial ones (those that can be performed immediately). At each iteration, identify predecessors and successors, and set the `minCount` of a process equal to the `minCount` of its successor times the quantity that its successor needs of the item by which they're linked, divided by the quantity that it produces of that item.

Now, define `maxCount` as the least common multiple of the denominators of all the `minCount`s, and set the `count` of a process equal to the numerator of `maxCount` times its `minCount`.

While the current task array is not empty, we check whether each task can be performed `count` number of times, given the resources. If not, it can't be scheduled any more. If so, it will be scheduled as soon as possible, given the durations and start times of its precursors. We consider initial processes first, then their successors, and so on, till the final process has been scheduled. Then we return to the beginning (the next pass), and keep going till there are no more resources to proceed.

For finite projects, the end time is defined as the start time of the final process plus its duration. For cycling projects, the provisional end time is updated as tasks are scheduled, and returned along with the schedule when the timer signals to finish.

The examples show that more than one instance of a process can be scheduled simultaneously, which makes time optimization trivial: just schedule as many instances of all tasks, in the necessary proportions, as resources permit.

## 6. Further

While this program does generate plausible schedules for the given examples and our own simple configuration files, it's far from robust. It doesn't yet allow for the possibility of one task having multiple successors. It assumes tasks have been well chosen and just need giving start times and number of instances to perform at those times. It doesn't decide effectively between rival processes, having the same input and output, whether of the same or differing effectiveness:

```board:7

do_doorknobs:(board:1):(doorknobs:1):15
do_more_doorknobs:(board:1):(doorknobs:2):15
do_background:(board:2):(background:1):20
do_shelf:(board:1):(shelf:1):10
do_cabinet:(doorknobs:2;background:1;shelf:3):(cabinet:1):30

optimize:(time;cabinet)
```

The above configuration results in too many doorknobs and no cabinet! We should really just pick the most effective of such rivals, but what if the effectiveness is only demonstrated several steps down the line? A more thorough version would also want to deal with mischievous processes, such as `do_nothing:(board:1):(cabinet:0):15` or `do_what_now?:(board:0):(caperberries:12):15`.

Note the sensitivity to task-listing order of our own example `macguffin`.

If we content ourselves with a heuristic, should we favor not wasting surplus when resources are plentiful, or should we make sure we obtain at least one unit of the end product however sparse they are?

A better scheme might be to consider the linking item and divide the quantity required by the quantity produced, then take the ceiling to obtain the minimum number of times the producer needs to be performed to contribute one unit of its successor, provided other requirements are met. If there are multiple linking items, as in `fertilizer`, we'd chose the maximum of these ceilings. Having found how many times each task needs to be executed to obtain a unit of the goal, a first pass of scheduling could be performed, and the resources updated. In this way, a solution could be found incrementally ...
