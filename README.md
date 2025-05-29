# stock-exchange-sim

1. [Brief](#0-brief)
2. [Setup](#1-setup)
3. [Usage](#2-usage)
4. [Audit](#3-audit)
5. [Research](#4-research)
6. [Strategy](#5-strategy)
7. [Further](#6-further)

## 1. Brief

This is a project in the [01 Edu](https://01-edu.org/) system, introducing the idea of a [process chain](https://en.wikipedia.org/wiki/Event-driven_process_chain). It's an optional extra, at least for us at 01Founders in London, 2024. It can be done in any compiled language. I chose Go.

According to the [instructions](https://github.com/01-edu/public/tree/master/subjects/stock-exchange-sim), we need to write a program, `stock`, that takes one or two command-line arguments.

The first argument is required. It represents the path to a text file, defining a project, such as building a cabinet (`examples/build`) or growing apples (`examples/fertilizer`), that can be broken down into smaller tasks.<sup id="ref-f1">[1](#f1)</sup> There are two types of project: those that have a finite objective, such as building the cabinet, and those that can continue indefinitely thanks to renewable resources, such as growing the apples, eating them, and planting the seeds to produce more apple trees.

A correctly formatted file will contain a list of resources, together with the quantities of each resource available initially. It will list processes (intermediate tasks) with information about which resources each process consumes, how many units of each resource the process consumes, what items it can produce (given those resources), and how many of each it results in. The file should also say how long each process takes. Finally, it will specify a goal, in the form of an item to be "optimized", i.e. maximized.

This 'configuration file', as the instructions call it, may specify that time should be optimized too. In the case of non-renewable resources, we take this to mean that time should be minimized provided that the maximum amount of the goal is produced. No definition is given for what this might mean in the case of cyclic<sup id="ref-f2">[2](#f2)</sup> projects. (How would one minimize the duration of a neverending project?) The only example of such a project does not mention time.

The comment below, found in some of the examples, suggests the possibility of multiple resources to optimize, but none of the examples actually realize that possibility. The line of the examples that cites the goal always has the format `optimize:(<stock_name>)` or `optimize:(time;<stock_name>)`. Although the instructions speak of "elements" to optimize, the format they specify is `optimize:(<stock_name>|time)`. No indication is given of how one would decide between conflicting goals. We could show precedence by the order they're listed in, but, for now, I've taken the easier path of assuming only one stock item is to be maximized.

```
# optimize time for 0 stock and no process possible,
# or maximize some products over a long delay
# optimize:(stock1;stock2;...)
```

In principle, the comment might be taken to indicate a choice between optimizing "time" to exhausting resources (regardless of what's produced) or optimizing "some products" (regardless of how long this takes), but no example includes time without also including a product to optimize. The former objective defies common sense; the latter is trivial, given that there is no conflict, so time can be optimized whether this is a stated requirement or not.

Our program should take an optional second argument, an integer representing the maximum number of seconds the program is to run for. Although syntactically optional, this argument is necessary, in practice, for projects with renewable resources that can be replenished indefinitely. If omitted, `stock` will run for 1s by default.

Given a configuration file entitled `examples/build` (assuming Unix-like path separators), our program, `stock`, should produce a text file `examples/build.log`, consisting of a schedule: a list of processes (possibly including several instances of the same process, possibly overlapping), the statement "No more process doable at", followed by an integer one unit greater than the duration of whole project, and a list of stock (resources and products) left at the end.

We're also asked to make a checker that will check the processes listed in a log file and confirm that there are enough resources to perform each task listed at its start time, as specified in the log file.

## 2. Setup

To build an executable file of the `stock` program, navigate into the `stock-exchange-sim` folder and run the command `go build -o stock`. You could also type `go run . <file>`, but I recommend building an executable first if you want to see that the timer is working. Otherwise you'd have to wait for the program to compile before it runs for the duration you specify.

## 3. Usage

Enter `./stock examples/simple` to create a schedule for the example called `simple`.<sup id="ref-f3">[3](#f3)</sup>

Run `./stock simple 10` to specify that the program should not take longer than 10 seconds.

I've chosen to implement the checker as part of the same program. To check `simple.log`, run `./stock -checker examples/simple examples/simple.log`.

## 4. Audit

A shell script, `audit.zsh`, is provided for your convenience to run the examples in the Functional section of the audit instructions. (If using Bash or another shell, you can change the file extension and first line, `#!/usr/bin/env zsh`, accordingly.) Alternatively, you're welcome to type any or all of the commands yourself. In either case, please check the resulting files and ask any questions you might have.

As mentioned, I implemented the checker as part of the main program. A boolean flag is used to select checker mode. See `main.go` for the code that deals with the flag and other arguments, and `checker.go` for the checker function itself.

You'll find the configuration files for the given examples in the `examples` folder, together with the two examples we were required to create: `zen` (finite) and `matryushka` (infinite). (I also made a bonus finite one: `macguffin`.)

The audit instructions provide some intentionally erroneous configuration files to assess whether my program handles errors as required. The given erroneous configuration files, to test error handling in the main `stock` program, are in `examples/errors`. The configuration file `testchecker` and the corresponding erroneous log file `testchecker.log`, to test error handling in checker mode, are in `examples/checkererror`.

The finite examples are all simple enough to confirm manually that the logs are correct. The infinite ones can be examined for a few complete runs through to get the idea.

Exact outputs may vary from those suggested in the audit questions, especially where time is not to be optimized, since, in that case, there is less constraint on how soon tasks can be scheduled. Thus, for `seller`, the audit suggestion takes a more leisurely approach, whereas my program schedules processes as soon as the precedence relations allow, because why not?

Please note that I interpreted the time parameter as marking when to end the schedule function itself. Writing the output file (and printing the result to the terminal, if you choose to uncomment those lines in `main.go`) happen after the schedule is made. This seemed like the most natural interpretation of the instructions, particularly as printing to the terminal is an optional extra for the convenience of viewing small outputs without having to open the log file. It has no bearing on the result of the audit, which just asks you to confirm that fewer processes are performed in 0.003s than 1s for the example `fertilizer`.

Note also that how much the program can accomplish in a given duration, will vary depending on your computer and what it's is doing in the background on a given occasion. On one occasion, you might find that `./stock /examples/matryushka 0.001` manages several hundred runs of its infinite project. On another, it might not have time for any.

The instructions say, "It is up to you to create and organize the display, it must nevertheless allow the understanding of the main actions carried out by the program." By default, I just write the result to a log file. If you'd like to see it in the console too, you can uncomment the indicated line in `main.go` (before building the binary). This works fine for naturally terminating projects, but be warned that, when let run for even a few miliseconds, the indefinitely looping ones, `fertilizer` and `matryushka`, can output hundreds of thousands of lines. So, if you're using the shell script, you might want to just view the log files.

In the section about the checker, the instructions say, "The display must indicate whether the sequence is correct, or indicate the cycle and the process which are causing the problem. In all cases, at the end of the program, stocks are displayed, as well as the last cycle."

But the examples shown, both in the instructions and the audit, are of the form

```Evaluating: 0:buy_materiel
Evaluating: 10:build_product
Evaluating: 40:delivery
Trace completed, no error detected.
```

(Thus displaying neither stock nor cycle.) I've chosen to follow the format exemplified rather than that described. Maybe the description was a mistake or meant to refer instead to the schedule-generator program, which does indeed list any remaining stock and the last cycle. (Or, rather, the one after the last, in accordance with the examples shown.)

## 5. Research

Following the recommendation of the project description, I consulted [PM Knowledge Center](https://www.pmknowledgecenter.com), a collection of resources on "Project Management and Dynamic Scheduling". I found further background reading necessary to fill in the gaps in the explanations there: in particular, [Kolisch (1994)](https://www.econstor.eu/bitstream/10419/155418/1/manuskript_344.pdf). These sources describe what's known as a Resource Constrained Project Scheduling Problem. The heuristic type of solution our instructions direct us towards is called Priority Rule Based Scheduling.

Before going into detail, I should note that the above sources (or Kolish, at least, who gives more detailed algorithms) assume that each task can only be performed once per project, whereas our program is expected to deal with cases where tasks can and should be performed more than once (in succession or simultaneously), if resorces allow, to optimize what needs optimizing. This meant that I couldn't directly apply either of the proposed scheduling methods. Nevertheless, in case it's of interest, and since these were the theoretical resources we were directed towards, here is a summary of the main techniques described at PM Knowledge Center.

In Priority Rule Based Scheduling, a graph of precedence relationships is drawn up: that is, a graph where tasks are nodes, and a directed edge from A to B means that commencement of B depends directly on completion of A. Activities are numbered in such a way that successors always have a greater activity number than their predecessors. A priority rule is chosen. Then a schedule is generated according to one of two schemes:

- Serial

- Parallel

A SERIAL schedule generation scheme with N tasks takes N steps. One task is chosen, at each step, from the set of available tasks and moved to the set of completed tasks. (A task is available if it's the direct the successor to a completed task, and current resources suffice to perform it.) If multiple tasks are available, one is chosen according to the priority rule. If several have equal priority, the one with the lowest activity number is selected.

A PARALLEL schedule generation scheme with N tasks takes at most N steps. At each step, we schedule zero or more activities. Tasks are partitioned into completed, in progress, and available. The schedule time associated with a step is calculated as the earliest completion time of the tasks that were in progress during the previous step. Tasks whose finish time is equal to the schedule time are moved from the set of tasks in progress to the set of completed tasks. This may make other tasks available. As long as tasks are available, they're chosen one by one, in order as in a serial scheme, and started at the current schedule time, then we move on to the next step. The algorithm terminates when all tasks are completed or in progress.

## 6. Strategy

After all that, neither scheme quite works for us, given the different underlying assumptions of our project: multiple instances of a task schedulable, possibly simultaneously. But we can take inspiration from them.

I started by making the simplifying assumption that, as in the examples we're given, processes can have multiple predecessors but only one sucessor. I also made some "good faith" assumptions about the configuration file, such as the absence of processes that don't contribute towards the goal. A more robust scheme would need to deal with such cases.

My program schedules tasks by taking as many passes through the precedence graph as resources permit. Before the first pass, an array (in Go terms, a slice) of the tasks currently being considered is initialized with tasks that can start immediately, in the first unit of time.

I define the `count` of a process as the number of times it needs to be scheduled to produce one unit of the target item. The `iterations` of a process will be the number of times it is eventually scheduled.

Set all `count`s to zero initially. Also give each process a field `minCount` that will be used to initialize the count. `minCount` will be of a home-made type, rational, representing a rational number. Set the `minCount` of each task to 1 initially. Then proceed backwards from final tasks (defined as those that directly produce the target item) to initial ones (those that can be performed immediately). At each iteration, identify predecessors and successors, and set the `minCount` of a process equal to the `minCount` of its successor times the quantity that its successor needs of the item by which they're linked, divided by the quantity that it produces of that item.

Now, define `maxCount` as the least common multiple of the denominators of all the `minCount`s, and set the `count` of a process equal to the numerator of `maxCount` times its `minCount`.

While the current task array is not empty, I check whether each task can be performed `count` number of times, given the resources. If not, it can't be scheduled any more. If so, it will be scheduled as soon as possible, given the durations and start times of its precursors. I consider initial processes first, then their successors, and so on, till the final process has been scheduled. Then I return to the beginning (the next pass), and keep going till there are no more resources to proceed.

For finite projects, the end time is defined as the start time of the final process plus its duration. For theoretically infinite projects, the provisional end time is updated as tasks are scheduled, and returned along with the schedule when the timer signals to finish.

The examples show that more than one instance of a process can be scheduled simultaneously, which makes time optimization trivial: just schedule as many instances of all tasks, in the necessary proportions, as resources permit.

## 7. Further

While this program does generate plausible schedules for the given examples and my own simple configuration files, it's far from robust. It doesn't allow for the possibility of one task having multiple successors. It assumes tasks have been well chosen and just need giving start times and number of instances to perform at those times. It doesn't decide effectively between rival processes having the same input and output, whether of the same or differing effectiveness:

```board:7

do_doorknobs:(board:1):(doorknobs:1):15
do_more_doorknobs:(board:1):(doorknobs:2):15
do_background:(board:2):(background:1):20
do_shelf:(board:1):(shelf:1):10
do_cabinet:(doorknobs:2;background:1;shelf:3):(cabinet:1):30

optimize:(time;cabinet)
```

The above configuration results in too many doorknobs and no cabinet! We should really just pick the most effective of such rivals, but what if the effectiveness is only demonstrated several steps down the line? A more thorough version would also want to deal with lazy processes, such as `do_nothing:(board:1):(cabinet:0):15`, or mischievous ones: `do_what_now?:(board:0):(caperberries:12):15`.

Note the sensitivity to task-listing order of my example `macguffin`.

It might be better to build a solution in a more incremental way: we could take each linking item and divide the quantity required by the quantity produced, then take the ceiling to obtain the minimum number of times the producer needs to be performed to allow one run of its successor, provided other requirements are met. If there are multiple linking items, as in `fertilizer`, we'd chose the maximum of these ceilings. Having found how many times each task needs to be executed to obtain a unit of the goal, a first pass of scheduling could be performed, and the resources updated.

In slightly more detail, given a chain of two tasks, if the output of the first is less than or equal to the input of the second, we could schedule the first as many times as it takes till it produces enough to perform task 2. Work back along the chain till some of the target item can be made. If several parallel chains are required to produce the target item, augment each in a loop till the resources are all used up, moving in order from the cheapest to the most costly chain. If there are intermediate intersections, where more than one chain is necessary, run the above scheduling procedure piecewise, working backwards from goal to initial resources.

This assumes an already chosen graph of tasks where intersections mark points where multiple tasks are necessary for a successor task to proceed. We could call such nodes "and" crossings. But, returning to the example at the start of this section, we'd ideally want a foolproof way to decide among collections of paths from initial resources to goal in cases where a task might get the same raw material from different predecessor tasks: "or" crossings. The two types of crossing needn't coincide. We'd want to choose the most productive path where there are choices, but only if the productive path didn't have side-effects that impacted performance elsewhere in the graph.

Each new complication could be explored by making simplifying assumptions about other aspects, such as whether we allow one task to comsume and produce many types of resource.

Finally, we could think how to deal with infinite project in a more general way. We could consider how to divide up an indefinite project into natural units, perhaps defined by milestones where proportions of resources repeat, and maximize production of the goal over such units.

<a id="f1" href="#ref-f1">1</a>: Let's use the word 'project' to mean the whole thing being scheduled, and 'task' or 'process' to mean one of the individual tasks that it's composed of. The instructions favor the term 'process'. I'll refer to the item to be maximized as the 'goal' or 'target item'.[↩](#ref-f1)

<a id="f2" href="#ref-f2">2</a>: Cyclic in the sense that they can repeat indefinitely. The instructions use the word 'cycle' to mean a unit of time, as specified in a configuration file. To avoid a clash, I'll use the term 'runs through' for discrete cycles/iterations of a (potentially neverending) project from initial resources to production of the item that is the goal, returning the system to its initial state, except possibly for an increased quantity of that target item. I'll reserve the term 'iteration' for instances of an individual process, such as `do_doorknobs`, being performed.[↩](#ref-f2)

<a id="f3" href="#ref-f3">3</a>: For simplicity, I'll assume from now on that you're on a Unix-like operating system. If not, e.g. on Windows, please make the necessary adjustments to file extensions and path dividers.[↩](#ref-f3)
