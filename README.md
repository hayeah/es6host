Chrome [Canary supports ES6 modules](https://medium.com/dev-channel/es6-modules-in-chrome-canary-m60-ba588dfb8ab7). It's pretty great. No more Rube Goldberg setup is necessary to develop modern JavaScript apps.

This is a simple Go web server to support common ES6 app structure.

# Install

```
go get -u github.com/hayeah/es6host
```

# Instruction

Include the `index.js` of your app in `index.html`:

```
<script type="module" src="index.js"></script>
```

Inside any source file, you may import a module using relative path:

```
import foo from "./foo";
console.log("hello ES6 module");
```

Run the web server:

```
es6host -root ./
```

# TODO

* Support for node_modules lookup.