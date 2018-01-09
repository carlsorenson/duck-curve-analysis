import nodeResolve from 'rollup-plugin-node-resolve';
import commonjs    from 'rollup-plugin-commonjs';
import uglify      from 'rollup-plugin-uglify';
import copy from 'rollup-plugin-copy';

export default {
  input: 'intermediate/main.js',
  output: { file: '../output/build.js', format: 'iife' }, // output a single application bundle
  sourceMap: false,
  onwarn: function(warning) {
    // Skip certain warnings

    // should intercept ... but doesn't in some rollup versions
    if ( warning.code === 'THIS_IS_UNDEFINED' ) { return; }

    // console.warn everything else
    console.warn( warning.message );
  },
  plugins: [
      nodeResolve({jsnext: true, module: true}),
      commonjs({
        include: 'node_modules/rxjs/**',
      }),
      uglify(),
      copy({
        "intermediate/index.html": "../output/index.html",
        "intermediate/index.html": "../output/explore/index.html",
        "intermediate/index.html": "../output/conclusions/index.html",
        "intermediate/index.html": "../output/warmup/index.html",
        "intermediate/styles.css": "../output/styles.css",
        "node_modules/core-js/client/shim.min.js": "../output/node_modules/core-js/client/shim.min.js",
        "node_modules/zone.js/dist/zone.js": "../output/node_modules/zone.js/dist/zone.js"
      })
  ]
};
