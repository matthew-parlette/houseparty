import os
import jira
import todoist

from contextlib import contextmanager
from rocketchat.api import RocketChatAPI

class Bot(object):
    class Config(object):
        def __init__(self, path="config"):
            self.path = path

        def get(self, key):
            if os.path.exists(os.path.join(self.path, key)):
                with open(os.path.join(self.path, key), 'r') as setting: return setting.read().strip()
            return None

    class Secret(Config):
        def __init__(self, path="secrets"):
            self.path = path

    def __init__(self, path='.'):
        """
        Load configuration and secrets
        """
        self.config = Bot.Config(path="{}/config".format(path))
        self.secrets = Bot.Secret(path="{}/secrets".format(path))

    @contextmanager
    def JiraAPI(self):
        options = {'server': self.config.get("jira-url")}
        api = None
        try:
            api = jira.JIRA(options, basic_auth=(
                self.secrets.get("jira-username"),
                self.secrets.get("jira-password")))
        except jira.exceptions.JIRAError:
            pass
        yield api

    @contextmanager
    def TodoistAPI(self):
        api = todoist.TodoistAPI(self.secrets.get("todoist-token"))
        yield api

    @contextmanager
    def TodoistProject(self, name):
        with self.TodoistAPI() as api:
            api.reset_state()
            project = None
            for p in api.sync()['projects']:
                if p['name'].lower() == name.lower():
                    project = p
                    break
            yield project

    @contextmanager
    def ChatAPI(self):
        api = RocketChatAPI(settings={
            "username": self.config.get("rocketchat-email"),
            "password": self.secrets.get("rocketchat-password"),
            "domain": self.config.get("rocketchat-url")})
        yield api

    @contextmanager
    def ChatRoom(self, name):
        with self.ChatAPI() as api:
            rooms = api.get_public_rooms() + api.get_private_rooms()
            room = None
            for chatroom in rooms:
                if chatroom["name"] == name:
                    room = api.get_room_info(chatroom["id"])
                    break
            yield room

    def message(self, message, rooms=[]):
        for name in rooms:
            with self.ChatRoom(name) as room:
                # contextmanager may return None if the room isn't found
                if room:
                    with self.ChatAPI() as api:
                        api.send_message(message=message, room_id=room["channel"]["_id"])
