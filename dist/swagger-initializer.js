window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container

  var openApiUrl = document.URL.substr(0,document.URL.indexOf("swagger-ui/")) + "openapi.yml";
  console.debug("fetching open api", openApiUrl);

  window.ui = SwaggerUIBundle({
    url: openApiUrl,
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
