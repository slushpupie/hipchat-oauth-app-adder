# hipchat-oauth-app-adder
HipChat addons generally need to be web-accessible by the HipChat server, but it does allow setting up a "private" addon behind firewalls.  The docs are skimpy on how to accomplish this, so this utility will assist with someone trying to do this. It can be useful for a few other situations too.

### Basic Flow

HipChat needs to get a [Capabilities Descriptor](https://www.hipchat.com/docs/apiv2/capabilities) document to describe the addon. Typically this is done at the install screen by providing the URL to the descriptor.  This URL may be a [data](https://en.wikipedia.org/wiki/Data_URI_scheme) URL where the content is the descriptor. One of the links in the document is the "install" URL, which will be called via the browser conducting the install.  This app will retrieve the OAuth client details and redirect the browser back to HipChat to complete the installation.

### How to use this

Fire up the app, it will create a web service on port 4000. Browse to http://localhost:4000/ and fill out the form. When you submit the form, it will display the JSON Capabilities Descriptor, the Base64-encoded version of the document, and a handy link to install the addon if you wish.

### Contributing

This was thrown together to meet a very specific need, and likely has numerous bugs. If you wish to improve this, just submit a Pull Request. 
