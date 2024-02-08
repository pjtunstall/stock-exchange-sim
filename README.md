TODO:

- Check logic.
- Refactor.
- Deal with the cyclic case (renewable resources), as in fertilizer.
  - This will need a somewhat different logic.
- Write checker.
- Create own finite example.
- Create own infinite example.

# stock-exchange-sim

0. [Brief](#-brief)
1. [Setup](#1-setup)
2. [Usage](#2-usage)
3. [Research](#3-research)
4. [Strategy](#4-strategy)

## 0. Brief

According to the [instructions](https://github.com/01-edu/public/tree/master/subjects/stock-exchange-sim), we need to write a program, `stock`, that takes two command-line arguments.

This first is the path to a text file, which they call a "configuration file", formatted in a certain way. Our program should derive information, from this file, about a project. The file will contain a list of resources and quantities of each available at the start of the project; a list of processes with information about which resources they consume, in what quantity, and what products they produce and in what quantity; also, how long each process takes. Finally, it will specify the goal of the project, in the form of an item whose production should be "optimized", i.e. maximized.

The configuration file can also specify that time should be optimized too. In the case of non-renewable resources, we take this to mean that time should be minimized given that the maximum amount of the goal should be produced.

Comments in some of the examples suggest the possibility of multiple resources to optimize, but none of the examples actually realize that possibility. The line of the examples that cites the goal always has the format `optimize:(<stock_name>)` or `optimize:(time;<stock_name>)`. Although the instructions speak of "elements" to optimize, the format they specify is `optimize:(<stock_name>|time)`. No indication is given of how one would decide between conflicting goals. We could show precedence by the order they're listed in, but, for now, have taken the easier path of assuming only one stock item is to be maximized.

The second argument is an integer representing the maximum number of seconds the program must take to execute.

There are two types of task: those that can continue indefinitely thanks to renewable resources, and those that have a finite objective. (Given that the number of instances of a task that can be scheduled at once is only limited by resources and precedence relations, time optimization would seem to only trivially relevant when resources are not renewable.)

Given a configuration file `build.txt`, our program, `stock`, should produce a text file `build.log`, consisting of a schedule: a list of processes (possibly including several instances of the same process, possibly overlapping), an integer one unit greater than the duration of whole project, and a list of resources and products left at the end.

## 1. Setup

## 2. Usage

Run `go . simple.txt` (or, after building an executable, `./stock simple.txt`) to create a schedule for the example called `simple.txt`.

Run `go . simple.txt 10` to specify that the program should not take longer than 10 seconds.

## 3. Research

As recommended, we consulted [PM Knowledge Center](https://www.pmknowledgecenter.com), a collection of resources on "Project Management and Dynamic Scheduling". We found further background reading necessary to fill in the gaps in the explanations there: in particular, [Kolisch (1994)](https://www.econstor.eu/bitstream/10419/155418/1/manuskript_344.pdf). These sources describe what's known as a Resource Constrained Project Scheduling Problem. The heuristic type of solution our instructions direct us towards is called Priority Rule Based Scheduling.

Before going into detail, we should note that the above sources (or Kolish, at least, who gives more detailed algorithms) assume that each task can only be performed once per project, whereas our program is expected to deal with cases where tasks can and should be performed more than once (in succession or simultaneously), if resorces allow, to optimize what needs optimizing.

In Priority Rule Based Scheduling, a graph of precedence relationships is drawn up: that is, a graph where tasks are nodes, and a directed edge from A to B means that commencement of B depends directly on completion of A. Activities are numbered in such a way that successors always have a greater activity number than their predecessors. A priority rule is chosen. Then a schedule is generated according to one of two schemes:

- Serial

- Parallel

A SERIAL schedule generation scheme with N tasks takes N steps. One task is chosen, at each step, from the set of available tasks and moved to the set of completed tasks. (A task is available if it is the direct the successor to a completed task, and current resources suffice to perform it.) If multiple tasks are available, one is chosen according to the priority rule. If several have equal priority, the one with the lowest activity number is selected.

A PARALLEL schedule generation scheme with N tasks takes at most N steps. At each step, we schedule zero or more activities. Tasks are partitioned into completed, in progress, and available. The schedule time associated with a step is calculated as the earliest completion time of the tasks that were in progress during the previous step. Tasks whose finish time is equal to the schedule time are moved from the set of tasks in progress to the set of completed tasks. This may make other tasks available. As long as tasks are available, they're chosen one by one, in order as in a serial scheme, and started at the current schedule time, then we move on to the next step. The algorithm terminates when all tasks are completed or in progress.

## 4. Strategy

After all that, neither scheme quite works for us, given the different underlying assumptions of our project: multiple instances of a task schedulable, possibly simultaneously. But we can take inspiration from them.

Let's start with the simplifying assumption that, as in our examples, tasks can have multiple predecessors but only one sucessor.

...

Let the count of a task be the number of times it's scheduled to be performed. Set all counts to zero initially. Also give each process a field minCount that will be used to initialize the count. minCount will be of a home-made type, rational, representing a rational number. Set minCount of a task to the product of the minCount's of all its direct and indirect successors, multiplied by the quantity its successor requires, and divided by the quantity of the item it produces that connects it to its successor. (This assumes that, if it's connected via multiple products, the ratios of how much is produced to how much needed are all the same.)

...

Make resources a map for ease of adding resources as they're created and deleting them as they're used up?

...

Define the strength of a task as count times the number of units of target product (goal) it leads to per unit of resource it consumes. Assuming only essential tasks make it into the graph, we could increment a task's count till it reaches the minimum needed to make one unit of the end product. Repeat till the resources are all used up. Then mark current tasks as completed, increment resources if they're renewable, then move on to the next step.

For example, in `build`, we'd set the count of doorknobs to 2, background to 1, and shelf to 3. Since 2 \* 1 + 1 \* 2 + 3 \* 1 = 7, and there are 7 boards, we stop there. Time doesn't matter as we can schedule tasks an unlimited amount of times at once, and the only resource is never replenished. We just need to note when the cabinet building starts and finishes.

...
