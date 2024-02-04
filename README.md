# stock-exchange-sim

0. [Brief](#-brief)
1. [Setup](#1-setup)
2. [Usage](#2-usage)
3. [Research](#3-research)

## 0. Brief

According to the [instructions](https://github.com/01-edu/public/tree/master/subjects/stock-exchange-sim), we need to write a program, `stock_exchange`, that takes two command-line arguments.

This first is the path to a text file, which they call a "configuration file", formatted in a certain way. Our program should derive information, from this file, about a task. There'll be a list of resources and quantities of each available at the start of the task. There'll be a list of processes with information about which resources they consume, in what quantity, and what products they produce and in what quantity. Also, how long each process takes.

A goal will be specified in the form of an item whose production should be optimized.

The configuration file can also specify that time should be optimized too. I interpret this to mean that, if time is specified in the goal, we should try to maximize units of the target item divided by units of time.

The second argument is an integer representing the maximum number of seconds the program must take to execute.

There are two types of task: those that can continue indefinitely, and those that have a finite objective.

## 1. Setup

## 2. Usage

## 3. Research

As recommended, we consulted PM Knowledge Center, a collection of resources on "Project Management and Dynamic Scheduling", in particular Mario Vanhoucke's 2012 article, [Optimizing regular scheduling objectives](https://www.pmknowledgecenter.com/node/256). This describes a heuristic approach consisting of three steps:

- Assess the project data.

- Choose a priority rule and rank activities accordingly.

- Generate a schedule.

To quote: "A schedule generation scheme makes use of the priority list constructed in the previous step and aims at the generation of a feasible schedule by extending the partial schedule (i.e. a schedule where only a subset of the activities has been assigned a starting and finishing time) in a stage-wise fashion. At the start of the heuristic scheduling process, the partial schedule is empty and all activities are available to be scheduled. Afterwards, activities are selected according to their priorities and are put in the schedule following the rules of the generation scheme."

Generation schemes are said to fall into two categories:

- Serial: selects the activities one by one from the list and schedules it as-soon-as-possible in the schedule.

- Parallel: selects at each predefined time period the activities available to be scheduled and schedules them in the list as long as enough resources are available.

Note that serial and parallel here refer to the process of generating a schedule from a priority list rather than the order of the activities themseves. Both generations schemes schedule tasks to happen in parallel.

Four types of project data (activity, network, scheduling, resource) each give rise to a variety of priority rules, from which an activity list can be constructed, leading to the goal. We'll list some possibilities here, but note that, in Mario's words, "for each activity list, the network structure has been taken into account. This means that an activity can never be placed in the list before another activity if the latter is a (direct or indirect) predecessor of the former." Thus, network structure takes precedence over the rules listed below.

Two rules based on activity data are

- Shortest Processing Time (SPT): Put the activities in an increasing order of their durations in the list.

- Longest Processing Time (LPT): Put the activities in a decreasing order of their durations in the list.

Network logic suggests rules such as

- Most Immediate Successors (MIS): Put the activities with the most direct successors first in the activity list.

- Most Total Successors (MTS): Put the activities with the most successors, direct or otherwise, first in the activity list.

- Least Non-Related Jobs (LNRJ): A job (or activity) is not related to another job if there is no precedence related path between the two activities in the project network. Activities are ranked in an increasing order of their number of non-related activities.

- Greatest Rank Positional Weight (GRPW): The GRPW is calculated as the sum of the duration of the activity and the durations of its immediate successors.

Scheduling data gives rise to

- Earliest Start Time (EST): Put the activities in an increasing order of their earliest start in the list.

- Earliest Finish Time (EFT): Put the activities in an increasing order of their earliest finish in the list.

- Latest Start Time (LST): Put the activities in an increasing order of their latest start in the list.

- Latest Finish Time (LFT): Put the activities in an increasing order of their latest finish in the list.

- Minimum Slack (MSLK): Put the activities in an increasing order of their slack value in the list.

Resource data:

- Greatest Resource Work Content (GRWC): Put the activities in a decreasing order of their work content in the list.

- Greatest Cumulative Resource Work Content (GCRWC): Put the activities in a decreasing order of the sum of their work content and the work content of all their immediate successors in the list.
