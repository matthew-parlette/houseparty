#!/usr/bin/env python3
import subprocess
from pprint import pformat

try:
    print("Installing pip modules...")
    result = subprocess.check_call(
        ['sudo', 'pip3', 'install', '-r', 'houseparty/requirements.txt'],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE)
except subprocess.CalledProcessError:
    subprocess.check_call(['sudo', 'pip3', 'install', '-r', 'houseparty/requirements.txt'])
print("Pip modules installed!")

from houseparty import Bot

bot = Bot(path='.')

print("Testing JIRA...")
with bot.JiraAPI() as api:
    projects = api.projects()
    assert projects
    print("Found {} projects".format(len(projects)))
    print("JIRA is working!")

print("Testing Todoist...")
with bot.TodoistProject("personal") as project:
    assert project
    print("Found project: {}".format(project["name"]))
    print("Todoist is working!")

print("Testing Rocket.Chat...")
with bot.ChatAPI() as api:
    me = api.get_my_info()
    assert me
    print("Successfully logged in as {}".format(me["name"]))

print("Testing message sending...")
bot.message(message="Test message from houseparty bot framework", rooms=['house-party'])
