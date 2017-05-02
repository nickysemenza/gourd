
const express = require('express');

const app = express();
const fs = require('fs');
const path = require('path');
const YAML = require('yamljs');
app.set('port', (process.env.PORT || 4000));

if (process.env.NODE_ENV === 'production') {
    app.use(express.static('build'));
}


app.get('/api/test', (req, res) => {
    res.json('hi');
});
app.get('/api/recipes', (req, res) => {
    fs.readdir('recipes',(err, files) => {
        let l = files.filter(file=>path.extname(file) == ".yaml").map(file => {
            return path.basename(file,'.yaml');
        });
        res.json(l);
    })
});
app.get('/api/recipes/:slug', (req, res) => {
    nativeObject = YAML.load(`recipes/${req.params.slug}.yaml`);
    res.json(nativeObject);
   // fs.readFile(`recipes/${req.params.slug}.yaml`, (err, fileData) => {
   //     if(err) {res.send(err); return}
   //     data = YAML.parse(fileData);
   //     res.json(data);
   // })
});


app.listen(app.get('port'), () => {
    console.log(`Find the server at: http://localhost:${app.get('port')}/`); // eslint-disable-line no-console
});