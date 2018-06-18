import os
import jira
import todoist

from contextlib import contextmanager

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
