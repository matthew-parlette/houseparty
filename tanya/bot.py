#!/usr/bin/env python3
import asyncio
import datetime
from houseparty import Bot

def priority(priority):
    """
    Convert a jira priority to a todoist priority
    """
    return {
        'Lowest': 3,
        'Low': 3,
        'Medium': 2,
        'High': 1,
        'Highest': 1,
    }.get(priority, 2)

async def run(bot, interval):
    while True:
        await create_tasks_from_jira(bot)
        print("Waiting {} seconds for the next run...".format(interval))
        await asyncio.sleep(interval)

async def create_tasks_from_jira(bot):
    with bot.JiraAPI() as jira:
        issues = jira.search_issues("assignee = currentUser() AND resolution = Unresolved")
        bot.message(
            message="I'm currently processing {} JIRA issues assigned to {}".format(
                len(issues), bot.config.get("jira-username")),
            rooms=["house-party"])
        for issue in issues:
            print("Processing issue {}...".format(issue.key))
            with bot.Task(
                subject="[{}] {} - {}".format(issue.key, issue.fields.summary, "parlette.us"),
                project_name=bot.config.get("todoist-project")) as task:
                # print("Processing task:\n{}".format(task))
                if task["due_date_utc"]:
                    due_utc = datetime.datetime.strptime(task["due_date_utc"], '%a %d %b %Y %H:%M:%S %z').replace(tzinfo=None)
                    now = datetime.datetime.utcnow()
                    if now > due_utc:
                        # Task is overdue, set it to be due today
                        task.update(date_string="tod")
                else:
                    # Due date is not set, set the due date to today
                    task.update(date_string="tod")
                # Update priority to match jira
                task.update(priority=priority(issue.fields.priority))
        bot.message(
            message="I've completed processing the {} JIRA issues assigned to {}".format(
                len(issues), bot.config.get("jira-username")),
            rooms=["house-party"])

bot = Bot(path=".")

bot.message(message="I'm online", rooms=["house-party"])

loop = asyncio.get_event_loop()
loop.run_until_complete(
    asyncio.gather(
        run(bot, float(bot.config.get("interval"))),
    )
)
loop.close()

bot.message(message="I'm going offline", rooms=["house-party"])