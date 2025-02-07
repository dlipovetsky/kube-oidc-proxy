local main = import './manifests/main.jsonnet';

function(cloud='google') main {
  cloud: cloud,
  // this will only run the google cluster
  clouds: {
    google: main.clouds.google,
  },
  base_domain: '.kubernetes.example.net',
  cert_manager+: {
    letsencrypt_contact_email:: 'certificates@example.net',
  },
  dex+: if $.master then {
    users: [
      $.dex.Password('admin@example.net', '$2y$10$i2.tSLkchjnpvnI73iSW/OPAVriV9BWbdfM6qemBM1buNRu81.ZG.'),  // plaintext: secure
    ],
    // This shows how to add dex connectors
    connectors: [
      $.dex.Connector('github', 'GitHub', 'github', {
        clientID: '0123',
        clientSecret: '4567',
        orgs: [{
          name: 'example-net',
        }],
      }),
    ],
  } else {
  },
}
