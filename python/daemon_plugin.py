# -*- coding: utf-8 -*-
from irc3.plugins.command import command
from irc3.plugins.cron import cron
import irc3
import re
import modules
import importlib
import markovify


@irc3.plugin
class Plugin(object):

    def __init__(self, bot):
        self.bot = bot
        self.health = 0
        self.monster_health = 0
        self.monster_attack = 0
        self.last_action = None
        self.wait = 0

    @irc3.event(irc3.rfc.JOIN)
    def say_hi(self, mask, channel, **kw):
        """Say hi when someone join a channel"""
        self.channel = channel
        if mask.nick != self.bot.nick:
            self.bot.privmsg(channel, 'Hi %s!' % mask.nick)
        else:
            self.bot.privmsg(channel, 'Hi!')

    @command(permission='view')
    def echo(self, mask, target, args):
        """Echo command

           %%echo <words>...
        """
        #print(args)
        yield ' '.join(args['<words>'])

    @irc3.extend
    def make_trump(self, *args):
        trump = open('trump.txt').read()
        model = markovify.NewlineText(trump)
        return model.make_sentence()

    @command(permission='view')
    def trump(self, mask, target, args):
        """Trump command

           %%trump
        """
        yield self.make_trump()

    @command(permission='view')
    def reload(self, mask, target, args):
        """reload command

           %%reload
        """
        importlib.invalidate_caches()
        importlib.reload(modules)
        importlib.invalidate_caches()
        modules.reload_()
        self.bot.reload()
        yield 'Modules have been reloaded'

    @cron('*/1 * * * *')
    def trump_cron(self):
        self.bot.privmsg('#loud_and_stupid', self.make_trump())

