import EmberRouter from '@ember/routing/router';
import config from 'corgi/config/environment';

export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  this.route('editor', function () {
    this.route('edit', { path: ':id' });
  });
  this.route('settings');
  this.route('register');
  this.route('login');

  this.route('articles', function () {
    this.route('article', { path: ':id' });
  });
  this.route('profile');
  this.route('error', { path: '/*path' });
});
