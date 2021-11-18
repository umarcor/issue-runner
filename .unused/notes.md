EXAMPLE: how to set an output value in a step and read it in the following

```yml
- run: echo ::set-output name=action_fruit::strawberry
   id: fruit
- uses: actions/hello-world-javascript-action@master
   with:
     who-to-greet: ${{ steps.fruit.outputs.action_fruit }}
```