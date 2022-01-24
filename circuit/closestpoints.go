package circuit

// ClosestPoints attempts to find the optimal circuit by:
// 1. (optional) Create the optimal convex perimeter.
// 2. Attach each non-perimeter point to its 2 closest points (this applies to all points if the perimeter was not created).
// 2a. If more points are within DistanceTo(farther point) + tspmodel.Threshold, attach them as well.
// 3. Group each set of attached points, to determine how many distinct circuits exist (e.g. if A<->B<->C<->D<->B would be one group of points, E<->F<->G<->E would be another)
// 3a. If the convex perimeter is created, treat that as its own group.
// 4. While more than one group exists, determine the most efficient way to combine two of the groups and merge them
// 4a. Find the point in each circuit closest to a point in the other circuit.
// 4b. Determine which of the edges attached to those points would produce the smallest distance increase when merged to the other circuit.
// 5. Clean up the resulting circuit:
// 5a. If a vertex is attached to more than 2 points, determine the egde among those points that would be most costly to remove the point from,
//     via the largest (closest to +infinity) distance increase/decrease, and detach it from the remaining points.
// 5b. If 5a would cause a point to be attached to only one or fewer points, attach it to its closest edge.
// 5c. If any edges intersect, find the combination of those 4 points that does not intersect and produces the shortest edges.
// 5d. (optional) Decide on clockwise or counter clockwise ordering; not required if convex perimeter has been created, since that will decide ordering.
type ClosestPoints struct {
}

// ClosestPoints attempts to find the optimal circuit by:
// 1. (optional) Create the optimal convex perimeter.
// 2. Attach each non-perimeter point to its 2 closest points that are not already attached to 2 points (this applies to all points if the perimeter was not created).
// 3. Group each set of attached points, to determine how many distinct circuits exist (e.g. if A<->B<->C<->D<->A would be one group of points, E<->F<->G<->E would be another)
// 3a. If the convex perimeter is created, treat that as its own group.
// 4. While more than one group exists, determine the most efficient way to combine two of the groups and merge them
// 4a. Find the point in each circuit closest to a point in the other circuit.
// 4b. Determine which of the edges attached to those points would produce the smallest distance increase when merged to the other circuit.
// 5. Clean up the resulting circuit:
// 5a. If any edges intersect, find the combination of those 4 points that does not intersect and produces the shortest edges.
// 5b. (optional) Decide on clockwise or counter clockwise ordering; not required if convex perimeter has been created, since that will decide ordering.
type ClosestDetachedPoints struct {
}
