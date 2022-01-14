import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import ENV from 'corgi/config/environment';

export default class LoginFormComponent extends Component {
  @service session;
  @service router;

  domains = ['cor.gi', 'ali.nk'];
  domainDefault = 'elga.io';

  @tracked payload = {
    title: '',
    domain: 'cor.gi',
    keyword: '',
    url: '',
  };

  @tracked status = {
    error: '',
    ok: '',
    message: '',
  };

  @tracked isCreating = null;

  @action
  async submit(e) {
    e.preventDefault();
    this.isCreating = true;

    sleep(5000).then(() => {
      console.log('2 isCreating');
      console.log(this.isCreating);
      this.isCreating = null;
    });
  }

  @action
  async submitE(e) {
    e.preventDefault();

    if (this.payload.domain === '') {
      this.payload.domain = this.domainDefault;
    }

    let createLink = await fetch(`${ENV.APP.apiHost}/api/v1/links/`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        title: this.payload.title,
        domain: this.payload.domain,
        keyword: this.payload.keyword,
        url: this.payload.url,
      }),
    });

    let response = await createLink;
    let data = await response.json();

    // Delete.
    console.log('data:');
    console.log(data);

    this.status.ok = response.ok;
    this.status.error = !response.ok;

    if (this.status.error) {
      this.status.message = data.message;
    }

    console.log('response:');
    console.log(response);

    console.log('this.status');
    console.log(this.status);
  }

  @action
  async getLinks() {
    let links = await fetch(`${ENV.APP.apiHost}/api/v1/links/`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    let response = await links;
    let data = await response.json();

    // Delete.
    console.log('data:');
    console.log(data);
  }
}

function sleep(time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}
