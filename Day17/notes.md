# Day 17 (3/6/2025)

## Convolution

- matrix mult
- break up matricies into rows
- give X rows to each node/process
- each node then pads a row above and below of just 0
- also pad 0 column to front and end

## Gradient Descent in Linear Regression

- optimization algorithm that tries to find max/min of objective function

### Steps

1. init random model params
2. compute gradient of cost function wrt each parameter
    - take partial derivatives of cost function wrt parameters
3. update params of model by taking steps in opposite direction of model
    - choose `hyperparameter learning rate` (alpha) which tells step size of gradient
4. repeat steps 2, 3

## Neural Networks / Deep Learning

- Input Layer
- Hidden Layer
- Output Layer

- input layer = each node has value -1 < x < 1
- hidden layer has weight for each node in the layer before it
  - multiply each previous nodes weight by that nodes value
- output layer also has weights for the last hidden layer

- forward pass can be represented as a matrix mult

### Training Loop

- forward pass
- calculate error
- back pass
- determine contribution of each node to error
- updating weights accordingly

how to determine error contribution

- take derivative of error wrt weight (Backpropagation)

how to minimize error

- Gradient Descent
