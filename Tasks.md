# TODO

## Generate initial model

- [X] Model of a point
- [X] Function for generating a random set of points
- [X] User input for number of points
- [X] Use the function and input to generate the initial data

## Determine center of points

- [X] Function for generating the midpoint of a set of points
- [X] Use the function to calculate and store the midpoint for subsequent operations

## Determine the farthest point from the midpoint

- [X] Function for determining the distance between two points
- [X] Function for determining the farthest point, in a set of points, from a point
- [X] Find and store the point farthest from the midpoint

## Determine the farthest point from the previous point

- [X] Find and store the point farthest from the point that is farthest from the midpoint

## Create line segments between the previous two points (a-to-b and b-to-a)

- [X] Model of a line segment
- [X] Model of a polygon
- [X] Construct line segment from a-to-b
- [X] Construct line segment from b-to-a
- [X] Initialize polygon with these line segments

## Prepare collections to track points' locations relative to the polygon

- [X] Collection of interior points
- [ ] Collection of exterior points
- [X] Collection of hull points

## Grow the polygon to its maximum convex hull

### Determine the distance from each external point to the polygon

- [X] Function for determining the distance from a point to a line segment
- [X] Model for tracking a point, line segment, and the distance from the point to the line segment
- [ ] Collection of distance models
- [ ] Find the edge closest to each point and add it to the collection

### Find the exterior point farthest from the polygon (ties do not matter)

- [ ] Function for determining the farthest point from a collection of distance models
- [ ] Find the point farthest from the polygon

### Add the point to the polygon, along its closest edge

- [X] Function to add a point to the polygon by splitting a specified edge
- [ ] Use the function to add the point to the polygon
- [ ] Remove the point from the collection of exterior points

### Check for exterior points that are now interior points

- [ ] Function to determine if a point is external to the polygon
- [ ] Move any exterior points that are now inside the polygon from the exterior collection to the interior collection

### Repeat until all points are either interior or hull points

- [ ] Repeat all steps in `Grow the polygon to its maximum convex hull` if there are any points in the exterior collection

## Track which points can be moved when shrinking the polygon

- [ ] Collection of fixed points (convex hull vertices)
- [X] Collection of internal points (points inside the polygon, not part of the perimeter)
- [X] Collection of maleable points (concave vertices)
- [X] Collections of internal and maleable points need to track the distance to the edge closest to each point
- [X] Collection of maleable points or edges that have been updated

## Determine each maleable point's initial distance from the polygon

- [X] Function to determine the distance increase resulting from adding a vertex to an edge
- [X] Function to determine the edge closest to a point
- [ ] Calculate the edge closest to each internal point, and update the collection

## Shrink the convex hull into the minimum concave polygon that can encompass the points

### Find the closest point to the polygon

- [ ] Find the internal point with the smallest distance increase
- [ ] Add the point to the polygon
- [ ] Move the point from the internal point collection to the maleable point collection

### Update the other maleable vertices

- [X] Function to check if an edge is closer to a point than its current closest edge, if the point is not part of the new edge
- [X] Function to move a point within a polygon, from one edge to another
- [ ] Move any maleable vertex closer to this edge than its current edge
- [ ] Track any maleable vertices moved in this manner
- [ ] For each moved vertex, repeat the steps in `Update the other maleable vertices`

### Update the internal vertices

- [ ] For each internal point update the distance of the closest edge by checking if any newly created or updated edges are closer

### Repeat until all internal points are vertices

- [ ] Repeat all steps in `Shrink the convex hull into the minimum concave polygon that can encompass the points` until there are no points in the internal point collection
